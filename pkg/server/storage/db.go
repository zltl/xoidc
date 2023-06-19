package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zltl/xoidc/pkg/snowflake"
)

type DB struct {
	db *sql.DB
}

func NewDB(url string) *DB {
	if url == "" {
		url = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			"postgres",
			"123456",
			"localhost",
			5432,
			"xoidc")
	}
	log.Tracef("openning db: %s", url)

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{db: db}
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	log.Debugf("query: %s args: %+v", query, args)
	rs, err := d.db.Query(query, args...)
	if err != nil {
		log.Errorf("query error: %s", err)
	}
	return rs, err
}

func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	log.Debugf("query row: %s args: %+v", query, args)
	return d.db.QueryRow(query, args...)
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	log.Debugf("exec: %s args: %+v", query, args)
	rs, err := d.db.Exec(query, args...)
	if err != nil {
		log.Errorf("exec error: %s", err)
	}
	return rs, err
}

func (d *DB) Prepare(query string) (*sql.Stmt, error) {
	log.Debugf("prepare: %s", query)
	stmt, err := d.db.Prepare(query)
	if err != nil {
		log.Errorf("prepare error: %s", err)
	}
	return stmt, err
}

func (d *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	log.Debugf("prepare context: %s", query)
	stmt, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		log.Errorf("prepare context error: %s", err)
	}
	return stmt, err
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	log.Debugf("query context: %s args: %+v", query, args)
	rs, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Errorf("query context error: %s", err)
	}
	return rs, err
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	log.Debugf("query row context: %s args: %+v", query, args)
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	log.Debugf("exec context: %s args: %+v", query, args)
	rs, err := d.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Errorf("exec context error: %s", err)
	}
	return rs, err
}

func (d *DB) GetUserByUsername(_username string) (*User, error) {
	var id int64
	var username, password, firstname, lastname, email, phone sql.NullString
	var emailVerified, phoneVerified, isAdmin sql.NullBool
	var createTime, updateTime sql.NullTime

	err := d.QueryRow(`SELECT 
	username, password, create_time, udpate_time, first_name, last_name, 
	email, phone, email_verified, phone_verified, is_admin
	FROM users WHERE username=$1`, username).
		Scan(&id, &username, &password, &createTime, &updateTime, &firstname, &lastname,
			&email, &phone, &emailVerified, &phoneVerified, &isAdmin)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &User{
		ID:            snowflake.ID(id).Base64Url(),
		Username:      username.String,
		Password:      password.String,
		FirstName:     firstname.String,
		LastName:      lastname.String,
		Email:         email.String,
		EmailVerified: emailVerified.Bool,
		Phone:         phone.String,
		PhoneVerified: phoneVerified.Bool,
		// PreferredLanguage: "",
		IsAdmin: isAdmin.Bool,
	}, nil
}

func (d *DB) GetUserByID(id string) (*User, error) {
	rid, err := snowflake.Parse(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var username, password, firstname, lastname, email, phone sql.NullString
	var emailVerified, phoneVerified, isAdmin sql.NullBool
	var createTime, updateTime sql.NullTime
	err = d.QueryRow(`SELECT
	username, password, create_time, udpate_time, first_name, last_name,
	email, phone, email_verified, phone_verified, is_admin
	FROM users WHERE id=$1`, rid.Int64()).Scan(&username, &password, &createTime, &updateTime, &firstname, &lastname,
		&email, &phone, &emailVerified, &phoneVerified, &isAdmin)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &User{
		ID:            snowflake.ID(rid.Int64()).Base64Url(),
		Username:      username.String,
		Password:      password.String,
		FirstName:     firstname.String,
		LastName:      lastname.String,
		Email:         email.String,
		EmailVerified: emailVerified.Bool,
		Phone:         phone.String,
		PhoneVerified: phoneVerified.Bool,
		// PreferredLanguage: "",
		IsAdmin: isAdmin.Bool,
	}, nil
}

func (d *DB) NewAuthRequest(ctx context.Context, authReq *AuthRequest) error {

	body, err := json.Marshal(authReq)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = d.ExecContext(ctx, `INSERT INTO auth_requests
	(id, creation_date, application_id, body, done)
	values($1, $2, $3, $4, $5)`,
		authReq.ID, authReq.CreationDate, authReq.ApplicationID, body, authReq.Done)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) GetAuthRequest(ctx context.Context, id string) (*AuthRequest, error) {
	var _id sql.NullString
	var creationDate sql.NullTime
	var done sql.NullBool
	var applicationID sql.NullString
	var body []byte

	err := d.QueryRowContext(ctx, `SELECT 
	id, creation_date, application_id, body, done
	FROM auth_requests WHERE id=$1`, id).
		Scan(&_id, &creationDate, &applicationID, &body, &done)
	if err != nil {
		return nil, err
	}

	authReq := &AuthRequest{}
	err = json.Unmarshal(body, authReq)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	authReq.ID = _id.String
	authReq.CreationDate = creationDate.Time
	authReq.ApplicationID = applicationID.String
	authReq.done = done.Bool

	return authReq, nil
}

func (d *DB) DeleteAuthRequest(ctx context.Context, id string) error {
	_, err := d.ExecContext(ctx, `DELETE FROM auth_requests WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) SetAuthRequestDone(ctx context.Context, id string, done bool) error {
	_, err := d.ExecContext(ctx, `UPDATE auth_requests SET done=$1 WHERE id=$2`, done, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) AddCode(ctx context.Context, code string, requestId string) error {
	_, err := d.ExecContext(ctx, `INSERT INTO codes
	(id, request_id, create_time)
	values($1, $2, now())`,
		code, requestId)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) GetRequestIdByCode(ctx context.Context, code string) (string, error) {
	err := d.QueryRowContext(ctx, `SELECT request_id FROM codes WHERE id=$1`, code).
		Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (d *DB) DeleteCodeByRequestId(ctx context.Context, requestId string) error {
	_, err := d.ExecContext(ctx, `DELETE FROM codes WHERE request_id=$1`, requestId)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) AddRefreshToken(ctx context.Context, refreshToken *RefreshToken, id string) error {
	amr, _ := json.Marshal(refreshToken.AMR)
	aud, _ := json.Marshal(refreshToken.Audience)
	scopes, _ := json.Marshal(refreshToken.Scopes)
	_, err := d.ExecContext(ctx, `INSERT INTO refresh_tokens
	(id, token, auth_time, amr, audience, user_id, application_id, expiration, scopes)
	values($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		id, refreshToken.Token, refreshToken.AuthTime, amr, aud,
		refreshToken.UserID, refreshToken.ApplicationID, refreshToken.Expiration,
		scopes)

	return err
}

func (d *DB) DeleteRefreshToken(ctx context.Context, id string) error {
	_, err := d.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteRefreshTokenByApplicationIDUserID(ctx context.Context, applicationID, userID string) error {
	_, err := d.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE application_id=$1 AND user_id=$2`,
		applicationID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) GetRefreshToken(ctx context.Context, id string) (*RefreshToken, error) {
	var amr, aud, scopes string
	var refreshToken RefreshToken
	err := d.QueryRowContext(ctx, `SELECT
	token, auth_time, amr, audience, user_id, application_id, expiration, scopes
	FROM refresh_tokens WHERE id=$1`, id).
		Scan(&refreshToken.Token, &refreshToken.AuthTime, &amr, &aud,
			&refreshToken.UserID, &refreshToken.ApplicationID, &refreshToken.Expiration, &scopes)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(amr), &refreshToken.AMR)
	json.Unmarshal([]byte(aud), &refreshToken.Audience)
	json.Unmarshal([]byte(scopes), &refreshToken.Scopes)
	return &refreshToken, nil
}

func (d *DB) DeleteToken(ctx context.Context, id string) error {
	_, err := d.ExecContext(ctx, `DELETE FROM tokens WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}
func (d *DB) DeleteTokenByRefreshTokenID(ctx context.Context, id string) error {
	_, err := d.ExecContext(ctx, `DELETE FROM tokens WHERE refresh_token_id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *DB) AddToken(ctx context.Context, id string, token *Token) error {
	auds, _ := json.Marshal(token.Audience)
	scopes, _ := json.Marshal(token.Scopes)
	_, err := s.ExecContext(ctx, `INSERT INTO tokens
	(id, application_id, subject, refresh_token_id, audience, expiration, scopes)
	values($1, $2, $3, $4, $5, $6, $7)`,
		id, token.ApplicationID, token.Subject, token.RefreshTokenID,
		string(auds), token.Expiration, string(scopes))
	if err != nil {
		return err
	}
	return nil
}

func (s *DB) GetToken(ctx context.Context, id string) (*Token, error) {
	token := &Token{}
	var auds, scopes string

	err := s.QueryRowContext(ctx, `SELECT application_id, subject, 
	refresh_token_id, audience, expiration, scopes
	FROM tokens WHERE id=$1`, id).
		Scan(&token.ApplicationID, &token.Subject, &token.RefreshTokenID,
			&auds, &token.Expiration, &scopes)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(auds), &token.Audience)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = json.Unmarshal([]byte(scopes), &token.Scopes)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return token, nil
}

func (s *DB) DeleteTokenBySubjectApplicationID(ctx context.Context, applicationID string, subject string) error {
	_, err := s.ExecContext(ctx, `DELETE FROM tokens WHERE subject=$1 AND application_id=$2`, subject, applicationID)
	if err != nil {
		return err
	}
	return nil
}

func (s *DB) GetClient(ctx context.Context, clientId string) (Client, error) {
	var client Client
	err := s.QueryRowContext(ctx, `SELECT secret, application_type, auth_method,
	access_token_type, dev_mode, id_token_userinfo_claims_assertion,
	clock_skew FROM clients WHERE id=$1`, clientId).Scan(
		&client.secret, &client.accessTokenType, &client.authMethod,
		&client.accessTokenType, &client.devMode, &client.idTokenUserinfoClaimsAssertion,
		&client.clockSkew,
	)
	if err != nil {
		return Client{}, err
	}
	// redirect_uris
	rowsuri, err := s.QueryContext(ctx, `SELECT redirect_uri FROM client_redirect_uris WHERE client_id=$1`, clientId)
	if err != nil {
		return Client{}, err
	}
	defer rowsuri.Close()
	for rowsuri.Next() {
		var uri string
		err := rowsuri.Scan(&uri)
		if err != nil {
			return Client{}, err
		}
		client.redirectURIs = append(client.redirectURIs, uri)
	}

	// response_types
	rowsresp, err := s.QueryContext(ctx, `SELECT response_type FROM client_response_types WHERE client_id=$1`, clientId)
	if err != nil {
		return Client{}, err
	}
	defer rowsresp.Close()
	for rowsresp.Next() {
		var respType string
		err := rowsresp.Scan(&respType)
		if err != nil {
			return Client{}, err
		}
		client.responseTypes = append(client.responseTypes, oidc.ResponseType(respType))
	}

	// grant_type
	rows, err := s.QueryContext(ctx, `SELECT grant_type FROM client_grant_types WHERE client_id=$1`, clientId)
	if err != nil {
		return Client{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var grantType string
		err := rows.Scan(&grantType)
		if err != nil {
			return Client{}, err
		}
		client.grantTypes = append(client.grantTypes, oidc.GrantType(grantType))
	}

	// 	postLogoutRedirectURIGlobs     []string
	rowslogout, err := s.QueryContext(ctx, `SELECT post_logout_redirect_uri_glob FROM 
	client_post_logout_redirect_uri_globs WHERE client_id=$1`, clientId)
	if err != nil {
		return Client{}, err
	}
	defer rowslogout.Close()
	for rowslogout.Next() {
		var logout string
		err := rowslogout.Scan(&logout)
		if err != nil {
			return Client{}, err
		}
		client.postLogoutRedirectURIGlobs = append(client.postLogoutRedirectURIGlobs, logout)
	}

	//	redirectURIGlobs               []string
	rowsredirect, err := s.QueryContext(ctx, `SELECT redirect_uri_glob FROM
	client_redirect_uri_globs WHERE client_id=$1`, clientId)
	if err != nil {
		return Client{}, err
	}
	defer rowsredirect.Close()
	for rowsredirect.Next() {
		var redirect string
		err := rowsredirect.Scan(&redirect)
		if err != nil {
			return Client{}, err
		}
		client.redirectURIGlobs = append(client.redirectURIGlobs, redirect)
	}

	return client, nil
}

func (s *DB) DeleteClient(ctx context.Context, clientId string) error {
	// clients
	s.ExecContext(ctx, `DELETE FROM clients WHERE id=$1`, clientId)
	// client_grant_types
	s.ExecContext(ctx, `DELETE FROM client_grant_types WHERE client_id=$1`, clientId)
	// client_redirect_uris
	s.ExecContext(ctx, `DELETE FROM client_redirect_uris WHERE client_id=$1`, clientId)
	// client_redirect_uri_globs
	s.ExecContext(ctx, `DELETE FROM client_redirect_uri_globs WHERE client_id=$1`, clientId)
	// client_response_types
	s.ExecContext(ctx, `DELETE FROM client_response_types WHERE client_id=$1`, clientId)
	// client_post_logout_redirect_uri_globs
	s.ExecContext(ctx, `DELETE FROM client_post_logout_redirect_uri_globs WHERE client_id=$1`, clientId)

	return nil
}

func (s *DB) UpdateClient(ctx context.Context, client *Client) error {
	// clients
	_, err := s.ExecContext(ctx, `UPDATE clients SET secret=$1, application_type=$2, auth_method=$3,
	access_token_type=$4, dev_mode=$5, id_token_userinfo_claims_assertion=$6, clock_skew=$7 WHERE id=$8`,
		client.secret, client.applicationType, client.authMethod,
		client.accessTokenType, client.devMode, client.idTokenUserinfoClaimsAssertion,
		client.clockSkew, client.id)
	if err != nil {
		return err
	}
	// client_grant_types
	s.ExecContext(ctx, `DELETE FROM client_grant_types WHERE client_id=$1`, client.id)
	for _, grantType := range client.grantTypes {
		_, err = s.ExecContext(ctx, `INSERT INTO client_grant_types
		(client_id, grant_type)
		values($1, $2)`,
			client.id, grantType)
		if err != nil {
			return err
		}
	}
	// client_redirect_uris
	s.ExecContext(ctx, `DELETE FROM client_redirect_uris WHERE client_id=$1`, client.id)
	for _, redirectURI := range client.redirectURIs {
		_, err = s.ExecContext(ctx, `INSERT INTO client_redirect_uris
		(client_id, redirect_uri)
		values($1, $2)`,
			client.id, redirectURI)
		if err != nil {
			return err
		}
	}
	// client_redirect_uri_globs
	s.ExecContext(ctx, `DELETE FROM client_redirect_uri_globs WHERE client_id=$1`, client.id)
	for _, redirectURIGlob := range client.redirectURIGlobs {
		_, err = s.ExecContext(ctx, `INSERT INTO client_redirect_uri_globs
		(client_id, redirect_uri_glob)
		values($1, $2)`,
			client.id, redirectURIGlob)
		if err != nil {
			return err
		}
	}
	// client_response_types
	s.ExecContext(ctx, `DELETE FROM client_response_types WHERE client_id=$1`, client.id)
	for _, responseType := range client.responseTypes {
		_, err = s.ExecContext(ctx, `INSERT INTO client_response_types
		(client_id, response_type)
		values($1, $2)`,
			client.id, responseType)
		if err != nil {
			return err
		}
	}
	// client_post_logout_redirect_uri_globs
	s.ExecContext(ctx, `DELETE FROM client_post_logout_redirect_uri_globs WHERE client_id=$1`, client.id)
	for _, postLogoutRedirectURIGlob := range client.postLogoutRedirectURIGlobs {
		_, err = s.ExecContext(ctx, `INSERT INTO client_post_logout_redirect_uri_globs
		(client_id, post_logout_redirect_uri_glob)
		values($1, $2)`,
			client.id, postLogoutRedirectURIGlob)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DB) AddClient(ctx context.Context, client *Client) error {
	// clients
	_, err := s.ExecContext(ctx, `INSERT INTO clients
	(id, secret, application_type, auth_method, access_token_type, dev_mode, id_token_userinfo_claims_assertion, clock_skew)
	values($1, $2, $3, $4, $5, $6, $7, $8)`,
		client.id, client.secret, client.applicationType, client.authMethod,
		client.accessTokenType, client.devMode, client.idTokenUserinfoClaimsAssertion,
		client.clockSkew)
	if err != nil {
		return err
	}
	// client_grant_types
	for _, grantType := range client.grantTypes {
		_, err = s.ExecContext(ctx, `INSERT INTO client_grant_types
		(client_id, grant_type)
		values($1, $2)`,
			client.id, grantType)
		if err != nil {
			return err
		}
	}
	// client_redirect_uris
	for _, redirectURI := range client.redirectURIs {
		_, err = s.ExecContext(ctx, `INSERT INTO client_redirect_uris
		(client_id, redirect_uri)
		values($1, $2)`,
			client.id, redirectURI)
		if err != nil {
			return err
		}
	}
	// client_redirect_uri_globs
	for _, redirectURIGlob := range client.redirectURIGlobs {
		_, err = s.ExecContext(ctx, `INSERT INTO client_redirect_uri_globs
		(client_id, redirect_uri_glob)
		values($1, $2)`,
			client.id, redirectURIGlob)
		if err != nil {
			return err
		}
	}
	// client_response_types
	for _, responseType := range client.responseTypes {
		_, err = s.ExecContext(ctx, `INSERT INTO client_response_types
		(client_id, response_type)
		values($1, $2)`,
			client.id, responseType)
		if err != nil {
			return err
		}
	}
	// client_post_logout_redirect_uri_globs
	for _, postLogoutRedirectURIGlob := range client.postLogoutRedirectURIGlobs {
		_, err = s.ExecContext(ctx, `INSERT INTO client_post_logout_redirect_uri_globs
		(client_id, post_logout_redirect_uri_glob)
		values($1, $2)`,
			client.id, postLogoutRedirectURIGlob)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DB) AddService(ctx context.Context, clientId string, keyid string, key string) error {
	_, err := s.ExecContext(ctx, `INSERT INTO services
	(client_id, key_id, key)
	values($1, $2, $3)`,
		clientId, keyid, key)
	if err != nil {
		return err
	}
	return nil
}

func (s *DB) GetService(ctx context.Context, clientId string, keyid string) (string, error) {
	var key string
	err := s.QueryRowContext(ctx, `SELECT key FROM services WHERE client_id=$1 AND key_id=$2`, clientId, keyid).Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *DB) AddUserCode(ctx context.Context, usercode, devicecode string) error {
	_, err := s.ExecContext(ctx, `INSERT INTO user_codes
	(user_code, device_code)
	values($1, $2)`,
		usercode, devicecode)
	return err
}

func (s *DB) GetUserCode(ctx context.Context, usercode string) (string, error) {
	var devicecode string
	err := s.QueryRowContext(ctx, `SELECT device_code FROM user_codes WHERE user_code=$1`, usercode).Scan(&devicecode)
	if err != nil {
		return "", err
	}
	return devicecode, nil
}

func (s *DB) AddDeviceCode(ctx context.Context, deviceCode string, d deviceAuthorizationEntry) error {
	stateStr, err := json.Marshal(d.state)
	if err != nil {
		log.Errorf("failed to marshal state: %v", err)
		return err
	}
	_, err = s.ExecContext(ctx, `INSERT INTO device_codes
	(device_code, user_code, state)
	values($1, $2, $3)`,
		deviceCode, d.userCode, string(stateStr))
	return err
}

func (s *DB) GetDeviceCode(ctx context.Context, deviceCode string) (deviceAuthorizationEntry, error) {
	var d deviceAuthorizationEntry
	var stateStr string
	err := s.QueryRowContext(ctx, `SELECT user_code, state FROM device_codes WHERE device_code=$1`,
		deviceCode).Scan(&d.userCode, &stateStr)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal([]byte(stateStr), &d.state)
	if err != nil {
		log.Errorf("failed to unmarshal state: %v", err)
		return d, err
	}
	return d, nil
}

func (s *DB) UpdateDeviceCode(ctx context.Context, deviceCOde string, d deviceAuthorizationEntry) error {
	stateStr, err := json.Marshal(d.state)
	if err != nil {
		log.Errorf("failed to marshal state: %v", err)
		return err
	}
	_, err = s.ExecContext(ctx, `UPDATE device_codes SET user_code=$1, state=$2 WHERE device_code=$3`,
		d.userCode, string(stateStr), deviceCOde)
	return err
}
