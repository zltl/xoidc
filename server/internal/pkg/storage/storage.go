package storage

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/google/uuid"
	"github.com/lib/pq"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
	log "github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/server/pkg/password"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

// serviceKey1 is a public key which will be used for the JWT Profile Authorization Grant
// the corresponding private key is in the service-key1.json (for demonstration purposes)
var serviceKey1 = &rsa.PublicKey{
	N: func() *big.Int {
		n, _ := new(big.Int).SetString("00f6d44fb5f34ac2033a75e73cb65ff24e6181edc58845e75a560ac21378284977bb055b1a75b714874e2a2641806205681c09abec76efd52cf40984edcf4c8ca09717355d11ac338f280d3e4c905b00543bdb8ee5a417496cb50cb0e29afc5a0d0471fd5a2fa625bd5281f61e6b02067d4fe7a5349eeae6d6a4300bcd86eef331", 16)
		return n
	}(),
	E: 65537,
}

// storage implements the op.Storage interface
// typically you would implement this as a layer on top of your database
// for simplicity this example keeps everything in-memory
type Storage struct {
	lock sync.Mutex
	// userStore     UserStore
	services      map[string]Service
	refreshTokens map[string]*RefreshToken
	signingKey    signingKey
	deviceCodes   map[string]deviceAuthorizationEntry
	userCodes     map[string]string
	serviceUsers  map[string]*Client

	PGHost     string
	PGPort     int
	PGUsername string
	PGPassword string
	PGDBName   string

	db *sql.DB
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

// open sql connection
func (s *Storage) Open() error {
	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.PGHost, s.PGPort, s.PGUsername, s.PGPassword, s.PGDBName)

	loggerdb := sqldblogger.OpenDriver(
		info,
		&pq.Driver{},
		logrusadapter.New(log.StandardLogger()),
	)

	s.db = loggerdb

	err := s.db.Ping()
	if err != nil {
		return err
	}

	s.refreshTokens = make(map[string]*RefreshToken)
	s.services = map[string]Service{
		"service": {
			keys: map[string]*rsa.PublicKey{
				"key1": serviceKey1,
			},
		},
	}
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	s.signingKey = signingKey{
		id:        uuid.NewString(),
		algorithm: jose.RS256,
		key:       key,
	}
	s.deviceCodes = make(map[string]deviceAuthorizationEntry)
	s.userCodes = make(map[string]string)
	s.serviceUsers = map[string]*Client{
		"sid1": {
			// id:     "sid1",
			secret: "verysecret",
			grantTypes: []oidc.GrantType{
				oidc.GrantTypeClientCredentials,
			},
			accessTokenType: op.AccessTokenTypeBearer,
		},
	}

	return nil
}

func (s *Storage) GetClient(ctx context.Context, id string) (*Client, error) {
	log.Tracef("GetClient: id=%s", id)

	clientID, err := uuid.Parse(id)
	if err != nil {
		log.Errorf("Parse: %v", err)
		return nil, err
	}
	return s.GetClientByUUID(ctx, clientID)
}

// CheckUsernamePassword implements the `authenticate` interface of the login
func (s *Storage) CheckUsernamePassword(username, passwordInput, reqid string) error {
	log.Tracef("CheckUsernamePassword: username=%s", username)

	requid, err := uuid.Parse(reqid)
	if err != nil {
		log.Error(err)
		return err
	}

	request, err := s.GetAuthRequestByUUID(context.TODO(), requid)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("request not found")
	}

	clientIDStr := request.GetClientID()
	clientID := uuid.MustParse(clientIDStr)

	us, err := s.GetUserByUsername(context.TODO(), username, clientID)
	if err != nil {
		log.Errorf("QueryPassword: %v", err)
		return err
	}
	passHash := us.Password
	match, err := password.ComparePasswordAndHash(passwordInput, passHash)
	if err != nil {
		log.Errorf("ComparePasswordAndHash: %v", err)
		return err
	}
	if match {
		request.UserID = us.ID
		request.IsDone = true

		err = s.UpdateAuthRequest(context.Background(), request)
		if err != nil {
			log.Errorf("UpdateAuthRequest: %v", err)
			return err
		}
		return nil
	}

	return fmt.Errorf("username or password wrong")
}

// CreateAuthRequest implements the op.Storage interface
// it will be called after parsing and validation of the authentication request
func (s *Storage) CreateAuthRequest(ctx context.Context, authReq *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	if len(authReq.Prompt) == 1 && authReq.Prompt[0] == "none" {
		// With prompt=none, there is no way for the user to log in
		// so return error right away.
		return nil, oidc.ErrLoginRequired()
	}

	if userID == "" {
		userID = "00000000-0000-0000-0000-000000000000"
	}

	log.Info("CreateAuthRequest, userID=", userID)
	// typically, you'll fill your storage / storage model with the information of the passed object
	request := authRequestToInternal(authReq, userID)

	log.Infof("request: %+v", request)
	rid, err := s.StoreAuthRequest(context.TODO(), request)
	if err != nil {
		log.Errorf("StoreAuthRequest: %v", err)
		return nil, err
	}
	request.ID = rid

	return request, nil
}

// AuthRequestByID implements the op.Storage interface
// it will be called after the Login UI redirects back to the OIDC endpoint
func (s *Storage) AuthRequestByUID(ctx context.Context, id uuid.UUID) (op.AuthRequest, error) {
	request, err := s.GetAuthRequestByUUID(ctx, id)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("request not found")
	}
	return request, nil
}

func (s *Storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	rid, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return s.AuthRequestByUID(ctx, rid)
}

// AuthRequestByCode implements the op.Storage interface
// it will be called after parsing and validation of the token request (in an authorization code flow)
func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	// for this example we read the id by code and then get the request by id
	requestID, err := func() (uuid.UUID, error) {
		s.lock.Lock()
		defer s.lock.Unlock()
		rid, err := s.CodeToRequestID(ctx, code)

		return rid, err
	}()
	if err != nil {
		return nil, fmt.Errorf("code invalid or expired")
	}
	return s.AuthRequestByUID(ctx, requestID)
}

// SaveAuthCode implements the op.Storage interface
// it will be called after the authentication has been successful and before redirecting the user agent to the redirect_uri
// (in an authorization code flow)
func (s *Storage) SaveAuthCode(ctx context.Context, id string, code string) error {
	// for this example we'll just save the authRequestID to the code
	rid, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.StoreCodeRequestID(ctx, code, rid)
	return err
}

// DeleteAuthRequest implements the op.Storage interface
// it will be called after creating the token response (id and access tokens) for a valid
// - authentication request (in an implicit flow)
// - token request (in an authorization code flow)
func (s *Storage) DeleteAuthRequest(ctx context.Context, id string) error {
	// you can simply delete all reference to the auth request
	s.lock.Lock()
	defer s.lock.Unlock()

	reqid, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.DeleteAuthRequestByUUID(ctx, reqid)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.DeleteCodeRequestIDByRequestID(ctx, reqid)
	if err != nil {
		log.Error(err)
	}

	return nil
}

// CreateAccessToken implements the op.Storage interface
// it will be called for all requests able to return an access token (Authorization Code Flow, Implicit Flow, JWT Profile, ...)
func (s *Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (string, time.Time, error) {
	var applicationID string
	switch req := request.(type) {
	case *AuthRequest:
		applicationID = req.GetClientID()
	case op.TokenExchangeRequest:
		applicationID = req.GetClientID()
	}

	token, err := s.accessToken(applicationID, "", request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		return "", time.Time{}, err
	}
	return token.ID.String(), token.Expiration, nil
}

// CreateAccessAndRefreshTokens implements the op.Storage interface
// it will be called for all requests able to return an access and refresh token (Authorization Code Flow, Refresh Token Request)
func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
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
		return accessToken.ID.String(), refreshToken, accessToken.Expiration, nil
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
	return accessToken.ID.String(), refreshToken, accessToken.Expiration, nil
}

func (s *Storage) exchangeRefreshToken(ctx context.Context, request op.TokenExchangeRequest) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
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

	return accessToken.ID.String(), refreshToken, accessToken.Expiration, nil
}

// TokenRequestByRefreshToken implements the op.Storage interface
// it will be called after parsing and validation of the refresh token request
func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	token, ok := s.refreshTokens[refreshToken]
	if !ok {
		return nil, fmt.Errorf("invalid refresh_token")
	}
	return RefreshTokenRequestFromBusiness(token), nil
}

// TerminateSession implements the op.Storage interface
// it will be called after the user signed out, therefore the access and refresh token of the user of this client must be removed
func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	clientid, err := uuid.Parse(clientID)
	if err != nil {
		log.Error(err)
		return err
	}
	userid, err := uuid.Parse(userID)
	if err != nil {
		log.Error(err)
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.DeleteRefreshTokenByApplicationAndSubject(ctx, tx, clientid, userid)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}
	err = s.DeleteTokenByApplicationAndSubject(ctx, tx, clientid, userid)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// GetRefreshTokenInfo looks up a refresh token and returns the token id and user id.
// If given something that is not a refresh token, it must return error.
func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	refreshToken, ok := s.refreshTokens[token]
	if !ok {
		return "", "", op.ErrInvalidRefreshToken
	}
	return refreshToken.UserID.String(), refreshToken.ID.String(), nil
}

// RevokeToken implements the op.Storage interface
// it will be called after parsing and validation of the token revocation request
func (s *Storage) RevokeToken(ctx context.Context, tokenIDOrToken string, userID string, clientID string) *oidc.Error {
	// a single token was requested to be removed
	tokenid, err := uuid.Parse(tokenIDOrToken)
	if err != nil {
		log.Error(err)
		return oidc.ErrServerError().WithDescription("could not parse token id")
	}

	clientid, err := uuid.Parse(clientID)
	if err != nil {
		log.Error(err)
		return oidc.ErrServerError().WithDescription("could not parse client id")
	}

	accessToken, err := s.QueryToken(ctx, tokenid)
	if err == nil {
		if accessToken.ApplicationID != clientid {
			return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
		}
		err = s.DeleteTokenByID(ctx, accessToken.ID)
		if err != nil {
			return oidc.ErrServerError().WithDescription("could not delete token")
		}
		return nil
	}
	refreshToken, err := s.QueryRefreshToken(ctx, tokenid)
	if err != nil {
		// if the token is neither an access nor a refresh token, just ignore it, the expected behaviour of
		// being not valid (anymore) is achieved
		return nil
	}
	if refreshToken.ApplicationID != clientid {
		return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
	}
	// if it is a refresh token, you will have to remove the access token as well
	_ = s.DeleteRefreshTokenByID(ctx, refreshToken.ID)
	_ = s.DeleteTokenByRefreshTokenID(ctx, refreshToken.ID)

	return nil
}

// SigningKey implements the op.Storage interface
// it will be called when creating the OpenID Provider
func (s *Storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
	// in this example the signing key is a static rsa.PrivateKey and the algorithm used is RS256
	// you would obviously have a more complex implementation and store / retrieve the key from your database as well
	return &s.signingKey, nil
}

// SignatureAlgorithms implements the op.Storage interface
// it will be called to get the sign
func (s *Storage) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	return []jose.SignatureAlgorithm{s.signingKey.algorithm}, nil
}

// KeySet implements the op.Storage interface
// it will be called to get the current (public) keys, among others for the keys_endpoint or for validating access_tokens on the userinfo_endpoint, ...
func (s *Storage) KeySet(ctx context.Context) ([]op.Key, error) {
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
	s.lock.Lock()
	defer s.lock.Unlock()
	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return RedirectGlobsClient(client), nil
}

// AuthorizeClientIDSecret implements the op.Storage interface
// it will be called for validating the client_id, client_secret on token or introspection requests
func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return err
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
	return nil
}

// SetUserinfoFromRequests implements the op.CanSetUserinfoFromRequest interface.  In the
// next major release, it will be required for op.Storage.
// It will be called for the creation of an id_token, so we'll just pass it to the private function without any further check
func (s *Storage) SetUserinfoFromRequest(ctx context.Context, userinfo *oidc.UserInfo, token op.IDTokenRequest, scopes []string) error {
	return s.setUserinfo(ctx, userinfo, token.GetSubject(), token.GetClientID(), scopes)
}

// SetUserinfoFromToken implements the op.Storage interface
// it will be called for the userinfo endpoint, so we read the token and pass the information from that to the private function
func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID, subject, origin string) error {
	tokenid, err := uuid.Parse(tokenID)
	if err != nil {
		return err
	}
	token, err := s.QueryToken(ctx, tokenid)
	if err != nil {
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
	return s.setUserinfo(ctx, userinfo, token.Subject.String(), token.ApplicationID.String(), token.Scopes)
}

// SetIntrospectionFromToken implements the op.Storage interface
// it will be called for the introspection endpoint, so we read the token and pass the information from that to the private function
func (s *Storage) SetIntrospectionFromToken(ctx context.Context, introspection *oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
	tokenid, err := uuid.Parse(tokenID)
	if err != nil {
		return err
	}

	token, err := s.QueryToken(ctx, tokenid)
	if err != nil {
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
			introspection.ClientID = token.ApplicationID.String()
			return nil
		}
	}
	return fmt.Errorf("token is not valid for this client")
}

// GetPrivateClaimsFromScopes implements the op.Storage interface
// it will be called for the creation of a JWT access token to assert claims for custom scopes
func (s *Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
	return s.getPrivateClaimsFromScopes(ctx, userID, clientID, scopes)
}

func (s *Storage) getPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
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
	s.lock.Lock()
	defer s.lock.Unlock()
	service, ok := s.services[clientID]
	if !ok {
		return nil, fmt.Errorf("clientID not found")
	}
	key, ok := service.keys[keyID]
	if !ok {
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
	return nil
}

// createRefreshToken will store a refresh_token in-memory based on the provided information
func (s *Storage) createRefreshToken(accessToken *Token, amr []string, authTime time.Time) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	token := &RefreshToken{
		ID:            accessToken.RefreshTokenID,
		Token:         accessToken.RefreshTokenID.String(),
		AuthTime:      authTime,
		AMR:           amr,
		ApplicationID: accessToken.ApplicationID,
		UserID:        accessToken.Subject,
		Audience:      accessToken.Audience,
		Expiration:    time.Now().Add(5 * time.Hour),
		Scopes:        accessToken.Scopes,
	}
	s.StoreRefreshToken(context.TODO(), token)
	return token.Token, nil
}

// renewRefreshToken checks the provided refresh_token and creates a new one based on the current
func (s *Storage) renewRefreshToken(currentRefreshToken string) (string, string, error) {
	curtokid, err := uuid.Parse(currentRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	refreshToken, err := s.QueryRefreshToken(context.TODO(), curtokid)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	// deletes the refresh token and all access tokens which were issued based on this refresh token
	_ = s.DeleteRefreshTokenByID(context.TODO(), curtokid)
	_ = s.DeleteTokenByRefreshTokenID(context.TODO(), curtokid)
	// creates a new refresh token based on the current one
	token := uuid.New()
	refreshToken.Token = token.String()
	refreshToken.ID = token

	s.StoreRefreshToken(context.TODO(), &refreshToken)
	return token.String(), refreshToken.ID.String(), nil
}

// accessToken will store an access_token in-memory based on the provided information
func (s *Storage) accessToken(applicationID, refreshTokenID, subject string, audience, scopes []string) (*Token, error) {
	apid, _ := uuid.Parse(applicationID)
	refid, _ := uuid.Parse(refreshTokenID)
	sub, _ := uuid.Parse(subject)

	token := &Token{
		ID:             uuid.New(),
		ApplicationID:  apid,
		RefreshTokenID: refid,
		Subject:        sub,
		Audience:       audience,
		Expiration:     time.Now().Add(5 * time.Minute),
		Scopes:         scopes,
	}
	s.SaveToken(context.Background(), token)

	return token, nil
}

// setUserinfo sets the info based on the user, scopes and if necessary the clientID
func (s *Storage) setUserinfo(ctx context.Context, userInfo *oidc.UserInfo, userID, clientID string, scopes []string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	user, err := s.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.Subject = user.ID.String()
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
	if request.GetRequestedTokenType() == "" {
		request.SetRequestedTokenType(oidc.RefreshTokenType)
	}

	// Just an example, some use cases might need this use case
	if request.GetExchangeSubjectTokenType() == oidc.IDTokenType && request.GetRequestedTokenType() == oidc.RefreshTokenType {
		return errors.New("exchanging id_token to refresh_token is not supported")
	}

	// Check impersonation permissions
	uid, err := uuid.Parse(request.GetExchangeSubject())
	if err != nil {
		return err
	}
	user, err := s.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}

	if request.GetExchangeActor() == "" && !user.IsAdmin {
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
	return nil
}

// GetPrivateClaimsFromScopesForTokenExchange implements the op.TokenExchangeStorage interface
// it will be called for the creation of an exchanged JWT access token to assert claims for custom scopes
// plus adding token exchange specific claims related to delegation or impersonation
func (s *Storage) GetPrivateClaimsFromTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) (claims map[string]interface{}, err error) {
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
	authReq, ok := req.(*AuthRequest) // Code Flow (with scope offline_access)
	if ok {
		return authReq.GetClientID(), authReq.AuthTime, authReq.GetAMR()
	}
	refreshReq, ok := req.(*RefreshTokenRequest) // Refresh Token Request
	if ok {
		return refreshReq.ApplicationID.String(), refreshReq.AuthTime, refreshReq.AMR
	}
	return "", time.Time{}, nil
}

// customClaim demonstrates how to return custom claims based on provided information
func customClaim(clientID string) map[string]interface{} {
	return map[string]interface{}{
		"client": clientID,
		"other":  "stuff",
	}
}

func appendClaim(claims map[string]interface{}, claim string, value interface{}) map[string]interface{} {
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
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := s.GetClient(ctx, clientID)
	if err != nil {
		return err
	}

	if _, ok := s.userCodes[userCode]; ok {
		return op.ErrDuplicateUserCode
	}

	s.deviceCodes[deviceCode] = deviceAuthorizationEntry{
		deviceCode: deviceCode,
		userCode:   userCode,
		state: &op.DeviceAuthorizationState{
			ClientID: clientID,
			Scopes:   scopes,
			Expires:  expires,
		},
	}

	s.userCodes[userCode] = deviceCode
	return nil
}

func (s *Storage) GetDeviceAuthorizatonState(ctx context.Context, clientID, deviceCode string) (*op.DeviceAuthorizationState, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	entry, ok := s.deviceCodes[deviceCode]
	if !ok || entry.state.ClientID != clientID {
		return nil, errors.New("device code not found for client") // is there a standard not found error in the framework?
	}

	return entry.state, nil
}

func (s *Storage) GetDeviceAuthorizationByUserCode(ctx context.Context, userCode string) (*op.DeviceAuthorizationState, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	entry, ok := s.deviceCodes[s.userCodes[userCode]]
	if !ok {
		return nil, errors.New("user code not found")
	}

	return entry.state, nil
}

func (s *Storage) CompleteDeviceAuthorization(ctx context.Context, userCode, subject string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	entry, ok := s.deviceCodes[s.userCodes[userCode]]
	if !ok {
		return errors.New("user code not found")
	}

	entry.state.Subject = subject
	entry.state.Done = true
	return nil
}

func (s *Storage) DenyDeviceAuthorization(ctx context.Context, userCode string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.deviceCodes[s.userCodes[userCode]].state.Denied = true
	return nil
}

// AuthRequestDone is used by testing and is not required to implement op.Storage
func (s *Storage) AuthRequestDone(id string) error {
	reqid, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
		return err
	}

	ctx := context.TODO()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	req, err := s.TXGetAuthRequestByUUID(ctx, tx, reqid)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return errors.New("request not found")
	}
	req.IsDone = true
	err = s.TXUpdateAuthRequest(context.TODO(), tx, req)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	_ = tx.Commit()
	return nil
}

func (s *Storage) ClientCredentials(ctx context.Context, clientID, clientSecret string) (op.Client, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	client, ok := s.serviceUsers[clientID]
	if !ok {
		return nil, errors.New("wrong service user or password")
	}
	if client.secret != clientSecret {
		return nil, errors.New("wrong service user or password")
	}

	return client, nil
}

func (s *Storage) ClientCredentialsTokenRequest(ctx context.Context, clientID string, scopes []string) (op.TokenRequest, error) {
	client, ok := s.serviceUsers[clientID]
	if !ok {
		return nil, errors.New("wrong service user or password")
	}

	return &oidc.JWTTokenRequest{
		Subject:  client.id.String(),
		Audience: []string{clientID},
		Scopes:   scopes,
	}, nil
}
