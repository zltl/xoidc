package storage

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/square/go-jose.v2"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
)

var (
	_ op.Storage                  = &Storage{}
	_ op.ClientCredentialsStorage = &Storage{}
)

// storage implements the op.Storage interface
// typically you would implement this as a layer on top of your database
// for simplicity this example keeps everything in-memory
type Storage struct {
	signingKey signingKey

	db *DB
}

type signingKey struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       *rsa.PrivateKey
}

func (s *signingKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *signingKey) Key() interface{} {
	return s.key
}

func (s *signingKey) ID() string {
	return s.id
}

type publicKey struct {
	signingKey
}

func (s *publicKey) ID() string {
	return s.id
}

func (s *publicKey) Algorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *publicKey) Use() string {
	return "sig"
}

func (s *publicKey) Key() interface{} {
	return &s.key.PublicKey
}

func NewStorage(userStore UserStore) *Storage {
	// TODO: load signing key from file/DB
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	return &Storage{
		signingKey: signingKey{
			id:        uuid.NewString(),
			algorithm: jose.RS256,
			key:       key,
		},
		db: NewDB(""),
		// TODO: new DB
	}
}

// CheckUsernamePassword implements the `authenticate` interface of the login
func (s *Storage) CheckUsernamePassword(username, password, id string) error {
	logrus.Tracef("CheckUsernamePassword: %s", username)

	request, err := s.db.GetAuthRequest(context.TODO(), id)
	if err != nil {
		logrus.Debugf("request not found: %v", err)
		return fmt.Errorf("request not found")
	}

	user, err := s.db.GetUserByUsername(username)
	if err != nil {
		logrus.Debugf("user not found: %v", err)
		return fmt.Errorf("username or password wrong")
	}
	logrus.Debugf("user: %v", user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logrus.Debugf("password wrong: %v", err)
		return fmt.Errorf("username or password wrong")
	}

	request.UserID = user.ID
	request.done = true
	return nil
}

func (s *Storage) CheckUsernamePasswordSimple(username, password string) error {
	logrus.Tracef("CheckUsernamePasswordSimple: %s", username)

	user, err := s.db.GetUserByUsername(username)
	if err != nil {
		logrus.Debugf("user not found: %v", err)
		return fmt.Errorf("username or password wrong")
	}
	logrus.Debugf("user: %v", user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logrus.Debugf("password wrong: %v", err)
		return fmt.Errorf("username or password wrong")
	}

	return nil
}

// CreateAuthRequest implements the op.Storage interface
// it will be called after parsing and validation of the authentication request
func (s *Storage) CreateAuthRequest(ctx context.Context, authReq *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	logrus.Tracef("CreateAuthRequest: %s, %+v", userID, *authReq)

	if len(authReq.Prompt) == 1 && authReq.Prompt[0] == "none" {
		// With prompt=none, there is no way for the user to log in
		// so return error right away.
		return nil, oidc.ErrLoginRequired()
	}

	// typically, you'll fill your storage / storage model with the information of the passed object
	request := authRequestToInternal(authReq, userID)

	// you'll also have to create a unique id for the request (this might be done by your database; we'll use a uuid)
	request.ID = uuid.NewString()

	// and save it in your database (for demonstration purposed we will use a simple map)
	err := s.db.NewAuthRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// finally, return the request (which implements the AuthRequest interface of the OP
	return request, nil
}

// AuthRequestByID implements the op.Storage interface
// it will be called after the Login UI redirects back to the OIDC endpoint
func (s *Storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	logrus.Tracef("AuthRequestByID: %s", id)

	request, err := s.db.GetAuthRequest(ctx, id)
	if err != nil {
		logrus.Debugf("request not found: %v", err)
		return nil, fmt.Errorf("request not found")
	}
	return request, nil
}

// AuthRequestByCode implements the op.Storage interface
// it will be called after parsing and validation of the token request (in an authorization code flow)
func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	logrus.Tracef("AuthRequestByCode: %s", code)
	requestID, err := s.db.GetRequestIdByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code invalid or expired")
	}
	return s.AuthRequestByID(ctx, requestID)
}

// SaveAuthCode implements the op.Storage interface
// it will be called after the authentication has been successful and before redirecting the user agent to the redirect_uri
// (in an authorization code flow)
func (s *Storage) SaveAuthCode(ctx context.Context, id string, code string) error {
	// for this example we'll just save the authRequestID to the code
	logrus.Tracef("SaveAuthCode: %s, %s", id, code)

	s.db.AddCode(ctx, code, id)
	return nil
}

// DeleteAuthRequest implements the op.Storage interface
// it will be called after creating the token response (id and access tokens) for a valid
// - authentication request (in an implicit flow)
// - token request (in an authorization code flow)
func (s *Storage) DeleteAuthRequest(ctx context.Context, id string) error {
	// you can simply delete all reference to the auth request
	logrus.Tracef("DeleteAuthRequest: %s", id)
	err := s.db.DeleteAuthRequest(ctx, id)
	if err != nil {
		logrus.Errorf("error deleting auth request: %v", err)
		return err
	}
	err = s.db.DeleteCodeByRequestId(ctx, id)
	if err != nil {
		logrus.Errorf("error deleting code: %v", err)
		return err
	}
	return nil
}

// CreateAccessToken implements the op.Storage interface
// it will be called for all requests able to return an access token (Authorization Code Flow, Implicit Flow, JWT Profile, ...)
func (s *Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (string, time.Time, error) {
	logrus.Tracef("CreateAccessToken: %+v", request)
	var applicationID string
	switch req := request.(type) {
	case *AuthRequest:
		// if authenticated for an app (auth code / implicit flow) we must save the client_id to the token
		applicationID = req.ApplicationID
	case op.TokenExchangeRequest:
		applicationID = req.GetClientID()
	}

	token, err := s.accessToken(applicationID, "", request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		return "", time.Time{}, err
	}
	return token.ID, token.Expiration, nil
}

// CreateAccessAndRefreshTokens implements the op.Storage interface
// it will be called for all requests able to return an access and refresh token (Authorization Code Flow, Refresh Token Request)
func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	logrus.Tracef("CreateAccessAndRefreshTokens: %+v", request)
	// generate tokens via token exchange flow if request is relevant
	if teReq, ok := request.(op.TokenExchangeRequest); ok {
		return s.exchangeRefreshToken(ctx, teReq)
	}

	// get the information depending on the request type / implementation
	applicationID, authTime, amr := getInfoFromRequest(request)

	// if currentRefreshToken is empty (Code Flow) we will have to create a new refresh token
	if currentRefreshToken == "" {
		refreshTokenID := uuid.NewString()
		accessToken, err := s.accessToken(applicationID, refreshTokenID, request.GetSubject(), request.GetAudience(), request.GetScopes())
		if err != nil {
			return "", "", time.Time{}, err
		}
		refreshToken, err := s.createRefreshToken(accessToken, amr, authTime)
		if err != nil {
			return "", "", time.Time{}, err
		}
		return accessToken.ID, refreshToken, accessToken.Expiration, nil
	}

	// if we get here, the currentRefreshToken was not empty, so the call is a refresh token request
	// we therefore will have to check the currentRefreshToken and renew the refresh token
	refreshToken, refreshTokenID, err := s.renewRefreshToken(currentRefreshToken)
	if err != nil {
		return "", "", time.Time{}, err
	}
	accessToken, err := s.accessToken(applicationID, refreshTokenID, request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		return "", "", time.Time{}, err
	}
	return accessToken.ID, refreshToken, accessToken.Expiration, nil
}

func (s *Storage) exchangeRefreshToken(ctx context.Context, request op.TokenExchangeRequest) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	logrus.Tracef("exchangeRefreshToken: %+v", request)
	applicationID := request.GetClientID()
	authTime := request.GetAuthTime()

	refreshTokenID := uuid.NewString()
	accessToken, err := s.accessToken(applicationID, refreshTokenID, request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshToken, err := s.createRefreshToken(accessToken, nil, authTime)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken.ID, refreshToken, accessToken.Expiration, nil
}

// TokenRequestByRefreshToken implements the op.Storage interface
// it will be called after parsing and validation of the refresh token request
func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	logrus.Tracef("TokenRequestByRefreshToken: %s", refreshToken)

	token, err := s.db.GetRefreshToken(ctx, refreshToken)
	// token, ok := s.refreshTokens[refreshToken]
	if err != nil {
		return nil, fmt.Errorf("invalid refresh_token")
	}
	return RefreshTokenRequestFromBusiness(token), nil
}

// TerminateSession implements the op.Storage interface
// it will be called after the user signed out, therefore the access and refresh token of the user of this client must be removed
func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	logrus.Tracef("TerminateSession: %s, %s", userID, clientID)
	s.db.DeleteTokenBySubjectApplicationID(ctx, clientID, userID)
	s.db.DeleteRefreshTokenByApplicationIDUserID(ctx, clientID, userID)
	return nil
}

// GetRefreshTokenInfo looks up a refresh token and returns the token id and user id.
// If given something that is not a refresh token, it must return error.
func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	logrus.Tracef("GetRefreshTokenInfo: %s, %s", clientID, token)
	refreshToken, err := s.db.GetRefreshToken(ctx, token)
	// refreshToken, ok := s.refreshTokens[token]
	if err != nil {
		return "", "", op.ErrInvalidRefreshToken
	}
	return refreshToken.UserID, refreshToken.ID, nil
}

// RevokeToken implements the op.Storage interface
// it will be called after parsing and validation of the token revocation request
func (s *Storage) RevokeToken(ctx context.Context, tokenIDOrToken string, userID string, clientID string) *oidc.Error {
	// a single token was requested to be removed
	logrus.Tracef("RevokeToken: %s, %s, %s", tokenIDOrToken, userID, clientID)

	accessToken, err := s.db.GetToken(ctx, tokenIDOrToken)
	// accessToken, ok := s.tokens[tokenIDOrToken] // tokenID
	if err == nil {
		if accessToken.ApplicationID != clientID {
			return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
		}
		// if it is an access token, just remove it
		// you could also remove the corresponding refresh token if really necessary
		s.db.DeleteToken(ctx, accessToken.ID)
		// delete(s.tokens, accessToken.ID)
		return nil
	}
	refreshToken, err := s.db.GetRefreshToken(ctx, tokenIDOrToken)
	// refreshToken, ok := s.refreshTokens[tokenIDOrToken] // token
	if err != nil {
		// if the token is neither an access nor a refresh token, just ignore it, the expected behaviour of
		// being not valid (anymore) is achieved
		return nil
	}
	if refreshToken.ApplicationID != clientID {
		return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
	}
	// if it is a refresh token, you will have to remove the access token as well
	s.db.DeleteRefreshToken(ctx, refreshToken.ID)
	// delete(s.refreshTokens, refreshToken.ID)
	s.db.DeleteTokenByRefreshTokenID(ctx, refreshToken.ID)
	return nil
}

// SigningKey implements the op.Storage interface
// it will be called when creating the OpenID Provider
func (s *Storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
	logrus.Tracef("SigningKey")
	// in this example the signing key is a static rsa.PrivateKey and the algorithm used is RS256
	// you would obviously have a more complex implementation and store / retrieve the key from your database as well
	return &s.signingKey, nil
}

// SignatureAlgorithms implements the op.Storage interface
// it will be called to get the sign
func (s *Storage) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	logrus.Tracef("SignatureAlgorithms")
	return []jose.SignatureAlgorithm{s.signingKey.algorithm}, nil
}

// KeySet implements the op.Storage interface
// it will be called to get the current (public) keys, among others for the keys_endpoint or for validating access_tokens on the userinfo_endpoint, ...
func (s *Storage) KeySet(ctx context.Context) ([]op.Key, error) {
	logrus.Tracef("KeySet")
	// as mentioned above, this example only has a single signing key without key rotation,
	// so it will directly use its public key
	//
	// when using key rotation you typically would store the public keys alongside the private keys in your database
	// and give both of them an expiration date, with the public key having a longer lifetime
	return []op.Key{&publicKey{s.signingKey}}, nil
}

// GetClientByClientID implements the op.Storage interface
// it will be called whenever information (type, redirect_uris, ...) about the client behind the client_id is needed
func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	logrus.Tracef("GetClientByClientID: %s", clientID)
	client, err := s.db.GetClient(ctx, clientID)
	// client, ok := s.clients[clientID]
	if err != nil {
		return nil, fmt.Errorf("client not found")
	}
	return RedirectGlobsClient(&client), nil
}

// AuthorizeClientIDSecret implements the op.Storage interface
// it will be called for validating the client_id, client_secret on token or introspection requests
func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	logrus.Tracef("AuthorizeClientIDSecret: %s, %s", clientID, clientSecret)
	client, err := s.db.GetClient(ctx, clientID)
	if err != nil {
		return fmt.Errorf("client not found")
	}
	// for this example we directly check the secret
	// obviously you would not have the secret in plain text, but rather hashed and salted (e.g. using bcrypt)
	if client.secret != clientSecret {
		return fmt.Errorf("invalid secret")
	}
	return nil
}

// SetUserinfoFromScopes implements the op.Storage interface.
// Provide an empty implementation and use SetUserinfoFromRequest instead.
func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID, clientID string, scopes []string) error {
	logrus.Tracef("SetUserinfoFromScopes: %s, %s, %s", userID, clientID, scopes)
	return nil
}

// SetUserinfoFromRequests implements the op.CanSetUserinfoFromRequest interface.  In the
// next major release, it will be required for op.Storage.
// It will be called for the creation of an id_token, so we'll just pass it to the private function without any further check
func (s *Storage) SetUserinfoFromRequest(ctx context.Context, userinfo *oidc.UserInfo, token op.IDTokenRequest, scopes []string) error {
	logrus.Tracef("SetUserinfoFromRequest: %s, %s", token.GetSubject(), scopes)
	return s.setUserinfo(ctx, userinfo, token.GetSubject(), token.GetClientID(), scopes)
}

// SetUserinfoFromToken implements the op.Storage interface
// it will be called for the userinfo endpoint, so we read the token and pass the information from that to the private function
func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID, subject, origin string) error {
	logrus.Tracef("SetUserinfoFromToken: %s, %s, %s", tokenID, subject, origin)
	token, ok := func() (*Token, bool) {
		token, err := s.db.GetToken(ctx, tokenID)
		if err != nil {
			return nil, false
		}
		return token, true
	}()
	if !ok {
		return fmt.Errorf("token is invalid or has expired")
	}
	// the userinfo endpoint should support CORS. If it's not possible to specify a specific origin in the CORS handler,
	// and you have to specify a wildcard (*) origin, then you could also check here if the origin which called the userinfo endpoint here directly
	// note that the origin can be empty (if called by a web client)
	//
	// if origin != "" {
	//	client, ok := s.clients[token.ApplicationID]
	//	if !ok {
	//		return fmt.Errorf("client not found")
	//	}
	//	if err := checkAllowedOrigins(client.allowedOrigins, origin); err != nil {
	//		return err
	//	}
	//}
	return s.setUserinfo(ctx, userinfo, token.Subject, token.ApplicationID, token.Scopes)
}

// SetIntrospectionFromToken implements the op.Storage interface
// it will be called for the introspection endpoint, so we read the token and pass the information from that to the private function
func (s *Storage) SetIntrospectionFromToken(ctx context.Context, introspection *oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
	logrus.Tracef("SetIntrospectionFromToken: %s, %s, %s", tokenID, subject, clientID)
	token, ok := func() (*Token, bool) {
		token, err := s.db.GetToken(ctx, tokenID)
		if err != nil {
			return nil, false
		}
		return token, true
	}()
	if !ok {
		return fmt.Errorf("token is invalid or has expired")
	}
	// check if the client is part of the requested audience
	for _, aud := range token.Audience {
		if aud == clientID {
			// the introspection response only has to return a boolean (active) if the token is active
			// this will automatically be done by the library if you don't return an error
			// you can also return further information about the user / associated token
			// e.g. the userinfo (equivalent to userinfo endpoint)

			userInfo := new(oidc.UserInfo)
			err := s.setUserinfo(ctx, userInfo, subject, clientID, token.Scopes)
			if err != nil {
				return err
			}
			introspection.SetUserInfo(userInfo)
			//...and also the requested scopes...
			introspection.Scope = token.Scopes
			//...and the client the token was issued to
			introspection.ClientID = token.ApplicationID
			return nil
		}
	}
	return fmt.Errorf("token is not valid for this client")
}

// GetPrivateClaimsFromScopes implements the op.Storage interface
// it will be called for the creation of a JWT access token to assert claims for custom scopes
func (s *Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
	logrus.Tracef("GetPrivateClaimsFromScopes: %s, %s, %s", userID, clientID, scopes)
	return s.getPrivateClaimsFromScopes(ctx, userID, clientID, scopes)
}

func (s *Storage) getPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
	logrus.Tracef("getPrivateClaimsFromScopes: %s, %s, %s", userID, clientID, scopes)
	for _, scope := range scopes {
		switch scope {
		case CustomScope:
			claims = appendClaim(claims, CustomClaim, customClaim(clientID))
		}
	}
	return claims, nil
}

// GetKeyByIDAndClientID implements the op.Storage interface
// it will be called to validate the signatures of a JWT (JWT Profile Grant and Authentication)
func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID, clientID string) (*jose.JSONWebKey, error) {
	logrus.Tracef("GetKeyByIDAndClientID: %s, %s", keyID, clientID)
	key, err := s.db.GetService(ctx, clientID, keyID)
	if err != nil {
		return nil, fmt.Errorf("key not found")
	}

	return &jose.JSONWebKey{
		KeyID: keyID,
		Use:   "sig",
		Key:   key,
	}, nil
}

// ValidateJWTProfileScopes implements the op.Storage interface
// it will be called to validate the scopes of a JWT Profile Authorization Grant request
func (s *Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	logrus.Tracef("ValidateJWTProfileScopes: %s, %s", userID, scopes)
	allowedScopes := make([]string, 0)
	for _, scope := range scopes {
		if scope == oidc.ScopeOpenID {
			allowedScopes = append(allowedScopes, scope)
		}
	}
	return allowedScopes, nil
}

// Health implements the op.Storage interface
func (s *Storage) Health(ctx context.Context) error {
	logrus.Tracef("Health")
	return nil
}

// createRefreshToken will store a refresh_token in-memory based on the provided information
func (s *Storage) createRefreshToken(accessToken *Token, amr []string, authTime time.Time) (string, error) {
	logrus.Tracef("createRefreshToken")
	token := &RefreshToken{
		ID:            accessToken.RefreshTokenID,
		Token:         accessToken.RefreshTokenID,
		AuthTime:      authTime,
		AMR:           amr,
		ApplicationID: accessToken.ApplicationID,
		UserID:        accessToken.Subject,
		Audience:      accessToken.Audience,
		Expiration:    time.Now().Add(5 * time.Hour),
		Scopes:        accessToken.Scopes,
	}
	s.db.AddRefreshToken(context.TODO(), token, token.ID)
	return token.Token, nil
}

// renewRefreshToken checks the provided refresh_token and creates a new one based on the current
func (s *Storage) renewRefreshToken(currentRefreshToken string) (string, string, error) {
	logrus.Tracef("renewRefreshToken %s", currentRefreshToken)
	refreshToken, err := s.db.GetRefreshToken(context.Background(), currentRefreshToken)
	// refreshToken, ok := s.refreshTokens[currentRefreshToken]
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	// deletes the refresh token and all access tokens which were issued based on this refresh token
	s.db.DeleteRefreshToken(context.Background(), currentRefreshToken)
	// delete(s.refreshTokens, currentRefreshToken)
	s.db.DeleteTokenByRefreshTokenID(context.TODO(), refreshToken.ID)
	// creates a new refresh token based on the current one
	token := uuid.NewString()
	refreshToken.Token = token
	refreshToken.ID = token
	s.db.AddRefreshToken(context.Background(), refreshToken, token)
	// s.refreshTokens[token] = refreshToken
	return token, refreshToken.ID, nil
}

// accessToken will store an access_token in-memory based on the provided information
func (s *Storage) accessToken(applicationID, refreshTokenID, subject string, audience, scopes []string) (*Token, error) {
	logrus.Tracef("accessToken %s, %s, %s, %s, %s", applicationID, refreshTokenID, subject, audience, scopes)
	token := &Token{
		ID:             uuid.NewString(),
		ApplicationID:  applicationID,
		RefreshTokenID: refreshTokenID,
		Subject:        subject,
		Audience:       audience,
		Expiration:     time.Now().Add(5 * time.Minute),
		Scopes:         scopes,
	}
	s.db.AddToken(context.Background(), token.ID, token)
	// s.tokens[token.ID] = token
	return token, nil
}

// setUserinfo sets the info based on the user, scopes and if necessary the clientID
func (s *Storage) setUserinfo(ctx context.Context, userInfo *oidc.UserInfo, userID, clientID string, scopes []string) (err error) {
	logrus.Tracef("setUserinfo: %+v, %s, %s, %s", userInfo, userID, clientID, scopes)
	user, err := s.db.GetUserByID(userID)
	// user := s.userStore.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.Subject = user.ID
		case oidc.ScopeEmail:
			userInfo.Email = user.Email
			userInfo.EmailVerified = oidc.Bool(user.EmailVerified)
		case oidc.ScopeProfile:
			userInfo.PreferredUsername = user.Username
			userInfo.Name = user.FirstName + " " + user.LastName
			userInfo.FamilyName = user.LastName
			userInfo.GivenName = user.FirstName
			userInfo.Locale = oidc.NewLocale(user.PreferredLanguage)
		case oidc.ScopePhone:
			userInfo.PhoneNumber = user.Phone
			userInfo.PhoneNumberVerified = user.PhoneVerified
		case CustomScope:
			// you can also have a custom scope and assert public or custom claims based on that
			userInfo.AppendClaims(CustomClaim, customClaim(clientID))
		}
	}
	return nil
}

// ValidateTokenExchangeRequest implements the op.TokenExchangeStorage interface
// it will be called to validate parsed Token Exchange Grant request
func (s *Storage) ValidateTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) error {
	logrus.Tracef("ValidateTokenExchangeRequest: %s", request)
	if request.GetRequestedTokenType() == "" {
		request.SetRequestedTokenType(oidc.RefreshTokenType)
	}

	// Just an example, some use cases might need this use case
	if request.GetExchangeSubjectTokenType() == oidc.IDTokenType && request.GetRequestedTokenType() == oidc.RefreshTokenType {
		return errors.New("exchanging id_token to refresh_token is not supported")
	}

	// Check impersonation permissions
	if request.GetExchangeActor() == "" {
		logrus.Tracef("GetExchangeActor nil")
		return errors.New("user doesn't have impersonation permission")
	}
	uinfo, err := s.db.GetUserByID(request.GetExchangeSubject())
	if err != nil || !uinfo.IsAdmin {
		logrus.Tracef("GetUserByID: uinfo: %+v", uinfo)
		return errors.New("user doesn't have impersonation permission")

	}

	allowedScopes := make([]string, 0)
	for _, scope := range request.GetScopes() {
		if scope == oidc.ScopeAddress {
			continue
		}

		if strings.HasPrefix(scope, CustomScopeImpersonatePrefix) {
			subject := strings.TrimPrefix(scope, CustomScopeImpersonatePrefix)
			request.SetSubject(subject)
		}

		allowedScopes = append(allowedScopes, scope)
	}

	request.SetCurrentScopes(allowedScopes)

	return nil
}

// ValidateTokenExchangeRequest implements the op.TokenExchangeStorage interface
// Common use case is to store request for audit purposes. For this example we skip the storing.
func (s *Storage) CreateTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) error {
	// TODO: store request for audit purposes
	logrus.Infof("TokenExchangeRequest: %+v", request)
	return nil
}

// GetPrivateClaimsFromScopesForTokenExchange implements the op.TokenExchangeStorage interface
// it will be called for the creation of an exchanged JWT access token to assert claims for custom scopes
// plus adding token exchange specific claims related to delegation or impersonation
func (s *Storage) GetPrivateClaimsFromTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) (claims map[string]interface{}, err error) {
	logrus.Tracef("GetPrivateClaimsFromTokenExchangeRequest: %s", request)
	claims, err = s.getPrivateClaimsFromScopes(ctx, "", request.GetClientID(), request.GetScopes())
	if err != nil {
		return nil, err
	}

	for k, v := range s.getTokenExchangeClaims(ctx, request) {
		claims = appendClaim(claims, k, v)
	}

	return claims, nil
}

// SetUserinfoFromScopesForTokenExchange implements the op.TokenExchangeStorage interface
// it will be called for the creation of an id_token - we are using the same private function as for other flows,
// plus adding token exchange specific claims related to delegation or impersonation
func (s *Storage) SetUserinfoFromTokenExchangeRequest(ctx context.Context, userinfo *oidc.UserInfo, request op.TokenExchangeRequest) error {
	logrus.Tracef("SetUserinfoFromTokenExchangeRequest: %s", request)
	err := s.setUserinfo(ctx, userinfo, request.GetSubject(), request.GetClientID(), request.GetScopes())
	if err != nil {
		return err
	}

	for k, v := range s.getTokenExchangeClaims(ctx, request) {
		userinfo.AppendClaims(k, v)
	}

	return nil
}

func (s *Storage) getTokenExchangeClaims(ctx context.Context, request op.TokenExchangeRequest) (claims map[string]interface{}) {
	logrus.Tracef("getTokenExchangeClaims: %s", request)
	for _, scope := range request.GetScopes() {
		switch {
		case strings.HasPrefix(scope, CustomScopeImpersonatePrefix) && request.GetExchangeActor() == "":
			// Set actor subject claim for impersonation flow
			claims = appendClaim(claims, "act", map[string]interface{}{
				"sub": request.GetExchangeSubject(),
			})
		}
	}

	// Set actor subject claim for delegation flow
	// if request.GetExchangeActor() != "" {
	// 	claims = appendClaim(claims, "act", map[string]interface{}{
	// 		"sub": request.GetExchangeActor(),
	// 	})
	// }

	return claims
}

// getInfoFromRequest returns the clientID, authTime and amr depending on the op.TokenRequest type / implementation
func getInfoFromRequest(req op.TokenRequest) (clientID string, authTime time.Time, amr []string) {
	logrus.Tracef("getInfoFromRequest: %s", req)
	authReq, ok := req.(*AuthRequest) // Code Flow (with scope offline_access)
	if ok {
		return authReq.ApplicationID, authReq.authTime, authReq.GetAMR()
	}
	refreshReq, ok := req.(*RefreshTokenRequest) // Refresh Token Request
	if ok {
		return refreshReq.ApplicationID, refreshReq.AuthTime, refreshReq.AMR
	}
	return "", time.Time{}, nil
}

// customClaim demonstrates how to return custom claims based on provided information
func customClaim(clientID string) map[string]interface{} {
	logrus.Tracef("customClaim: %s", clientID)
	return map[string]interface{}{
		"client": clientID,
		"other":  "stuff",
	}
}

func appendClaim(claims map[string]interface{}, claim string, value interface{}) map[string]interface{} {
	logrus.Tracef("appendClaim: %s", claim)
	if claims == nil {
		claims = make(map[string]interface{})
	}
	claims[claim] = value
	return claims
}

type deviceAuthorizationEntry struct {
	deviceCode string
	userCode   string
	state      *op.DeviceAuthorizationState
}

func (s *Storage) StoreDeviceAuthorization(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes []string) error {
	logrus.Tracef("StoreDeviceAuthorization: %s ", clientID)
	_, err := s.db.GetClient(ctx, clientID)
	if err != nil {
		return errors.New("client not found")
	}

	if _, err := s.db.GetUserCode(ctx, userCode); err != nil {
		return op.ErrDuplicateUserCode
	}

	// s.deviceCodes[deviceCode] =
	dc := deviceAuthorizationEntry{
		deviceCode: deviceCode,
		userCode:   userCode,
		state: &op.DeviceAuthorizationState{
			ClientID: clientID,
			Scopes:   scopes,
			Expires:  expires,
		},
	}
	s.db.AddDeviceCode(ctx, deviceCode, dc)
	s.db.AddUserCode(ctx, userCode, deviceCode)
	return nil
}

func (s *Storage) GetDeviceAuthorizatonState(ctx context.Context, clientID, deviceCode string) (*op.DeviceAuthorizationState, error) {
	logrus.Tracef("GetDeviceAuthorizatonState: %s ", clientID)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	entry, err := s.db.GetDeviceCode(ctx, deviceCode)
	// entry, ok := s.deviceCodes[deviceCode]
	if err != nil || entry.state.ClientID != clientID {
		return nil, errors.New("device code not found for client") // is there a standard not found error in the framework?
	}

	return entry.state, nil
}

func (s *Storage) GetDeviceAuthorizationByUserCode(ctx context.Context, userCode string) (*op.DeviceAuthorizationState, error) {
	logrus.Tracef("GetDeviceAuthorizationByUserCode: %s ", userCode)

	userCode, err := s.db.GetUserCode(ctx, userCode)
	if err != nil {
		return nil, errors.New("user code not found 1")
	}
	entry, err := s.db.GetDeviceCode(ctx, userCode)
	//entry, ok := s.deviceCodes[s.userCodes[userCode]]
	if err != nil {
		return nil, errors.New("user code not found 2")
	}

	return entry.state, nil
}

func (s *Storage) CompleteDeviceAuthorization(ctx context.Context, userCode, subject string) error {
	logrus.Tracef("CompleteDeviceAuthorization: %s ", userCode)

	userCode, err := s.db.GetUserCode(ctx, userCode)
	if err != nil {
		return errors.New("user code not found 3")
	}
	entry, err := s.db.GetDeviceCode(ctx, userCode)

	// entry, ok := s.deviceCodes[s.userCodes[userCode]]
	if err != nil {
		return errors.New("user code not found 4")
	}

	entry.state.Subject = subject
	entry.state.Done = true
	return nil
}

func (s *Storage) DenyDeviceAuthorization(ctx context.Context, userCode string) error {
	logrus.Tracef("DenyDeviceAuthorization: %s ", userCode)
	userCode, err := s.db.GetUserCode(ctx, userCode)
	if err != nil {
		return errors.New("user code not found 5")
	}
	entry, err := s.db.GetDeviceCode(ctx, userCode)
	if err != nil {
		return errors.New("user code not found 6")
	}
	entry.state.Denied = true

	s.db.UpdateDeviceCode(ctx, userCode, entry)
	// s.deviceCodes[s.userCodes[userCode]].state.Denied = true
	return nil
}

// AuthRequestDone is used by testing and is not required to implement op.Storage
func (s *Storage) AuthRequestDone(id string) error {
	logrus.Tracef("AuthRequestDone: %s ", id)
	_, err := s.db.GetAuthRequest(context.TODO(), id)
	if err != nil {
		s.db.SetAuthRequestDone(context.TODO(), id, true)
		// req.done = true
		return nil
	}

	return errors.New("request not found")
}

func (s *Storage) ClientCredentials(ctx context.Context, clientID, clientSecret string) (op.Client, error) {
	logrus.Tracef("ClientCredentials: %s ", clientID)

	client, err := s.db.GetClient(ctx, clientID)
	// client, ok := s.serviceUsers[clientID]
	if err != nil {
		return nil, errors.New("wrong service user or password")
	}
	if client.secret != clientSecret {
		return nil, errors.New("wrong service user or password")
	}

	return &client, nil
}

func (s *Storage) ClientCredentialsTokenRequest(ctx context.Context, clientID string, scopes []string) (op.TokenRequest, error) {
	logrus.Tracef("ClientCredentialsTokenRequest: %s ", clientID)
	client, err := s.db.GetClient(ctx, clientID)
	if err != nil {
		return nil, errors.New("wrong service user or password")
	}

	return &oidc.JWTTokenRequest{
		Subject:  client.id,
		Audience: []string{clientID},
		Scopes:   scopes,
	}, nil
}
