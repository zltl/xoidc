package storage

import (
	"context"
	"log"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sanyokbig/pqinterval"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var (
	// we use the default login UI and pass the (auth request) id
	defaultLoginURL = func(id string) string {
		return "/login/username?authRequestID=" + id
	}
)

func (s *Storage) TotalClient(ctx context.Context) (int64, error) {
	cmd := `
	SELECT
		count(*)
	FROM
		client
	`
	var total int64
	err := s.db.QueryRowContext(ctx, cmd).Scan(&total)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return total, nil
}

func (s *Storage) GetAllClient(ctx context.Context, offset, count int64) ([]Client, error) {
	cmd := `
	SELECT
		id,
		secret,
		redirect_uris,
		application_type,
		auth_method,
		response_types,
		grant_types,
		access_token_type,
		dev_mode,
		id_token_user_info_claims_assertion,
		clock_skew,
		post_logout_redirect_uri_globs,
		redirect_uri_globs,
		user_namespace_id,
		name
	FROM
		client
	LIMIT $1 OFFSET $2
	`
	rows, err := s.db.QueryContext(ctx, cmd, count, offset)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	var clients []Client

	for rows.Next() {
		c := &Client{}
		var interval pqinterval.Interval
		err := rows.Scan(
			&c.id,
			&c.secret,
			pq.Array(&c.redirectURIs),
			(*int)(unsafe.Pointer(&c.applicationType)),
			&c.authMethod,
			pq.Array((*[]string)(unsafe.Pointer(&c.responseTypes))),
			pq.Array((*[]string)(unsafe.Pointer(&c.grantTypes))),
			(*int)(unsafe.Pointer(&c.accessTokenType)),
			&c.devMode,
			&c.idTokenUserinfoClaimsAssertion,
			&interval,
			pq.Array(&c.postLogoutRedirectURIGlobs),
			pq.Array(&c.redirectURIGlobs),
			&c.userNamespaceID,
			&c.name,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		dura, err := interval.Duration()
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		c.clockSkew = dura
		c.loginURL = defaultLoginURL

		clients = append(clients, *c)
	}
	return clients, nil
}

func (s *Storage) GetClientByUUID(ctx context.Context, clientID uuid.UUID) (*Client, error) {
	stmt := `
		SELECT
			id,
			secret,
			redirect_uris,
			application_type,
			auth_method,
			response_types,
			grant_types,
			access_token_type,
			dev_mode,
			id_token_user_info_claims_assertion,
			clock_skew,
			post_logout_redirect_uri_globs,
			redirect_uri_globs,
			user_namespace_id,
			name
		FROM
			client
		WHERE
			id = $1
	`
	c := &Client{}
	var interval pqinterval.Interval
	err := s.db.QueryRowContext(ctx, stmt, clientID).Scan(
		&c.id,
		&c.secret,
		pq.Array(&c.redirectURIs),
		(*int)(unsafe.Pointer(&c.applicationType)),
		&c.authMethod,
		pq.Array((*[]string)(unsafe.Pointer(&c.responseTypes))),
		pq.Array((*[]string)(unsafe.Pointer(&c.grantTypes))),
		(*int)(unsafe.Pointer(&c.accessTokenType)),
		&c.devMode,
		&c.idTokenUserinfoClaimsAssertion,
		&interval,
		pq.Array(&c.postLogoutRedirectURIGlobs),
		pq.Array(&c.redirectURIGlobs),
		&c.userNamespaceID,
		&c.name,
	)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	dura, err := interval.Duration()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	c.clockSkew = dura

	c.loginURL = defaultLoginURL

	return c, err
}

// Client represents the storage model of an OAuth/OIDC client
// this could also be your database model
type Client struct {
	id                             uuid.UUID
	secret                         string
	redirectURIs                   []string
	applicationType                op.ApplicationType
	authMethod                     oidc.AuthMethod
	loginURL                       func(string) string
	responseTypes                  []oidc.ResponseType
	grantTypes                     []oidc.GrantType
	accessTokenType                op.AccessTokenType
	devMode                        bool
	idTokenUserinfoClaimsAssertion bool
	clockSkew                      time.Duration
	postLogoutRedirectURIGlobs     []string
	redirectURIGlobs               []string
	userNamespaceID                uuid.UUID
	name                           string
}

// RegisterClients enables you to register clients for the example implementation
// there are some clients (web and native) to try out different cases
// add more if necessary
//
// RegisterClients should be called before the Storage is used so that there are
// no race conditions.
func RegisterClients(registerClients ...*Client) {
	for _, client := range registerClients {
		log.Printf("Registering client: %s", client.id)
		// clients[client.id] = client
	}
}

// NativeClient will create a client of type native, which will always use PKCE and allow the use of refresh tokens
// user-defined redirectURIs may include:
// - http://localhost without port specification (e.g. http://localhost/auth/callback)
// - custom protocol (e.g. custom://auth/callback)
// (the examples will be used as default, if none is provided)
func NativeClient(id string, redirectURIs ...string) *Client {
	if len(redirectURIs) == 0 {
		redirectURIs = []string{
			"http://localhost/auth/callback",
			"custom://auth/callback",
		}
	}
	return &Client{
		id:                             uuid.MustParse(id),
		secret:                         "", // no secret needed (due to PKCE)
		redirectURIs:                   redirectURIs,
		applicationType:                op.ApplicationTypeNative,
		authMethod:                     oidc.AuthMethodNone,
		loginURL:                       defaultLoginURL,
		responseTypes:                  []oidc.ResponseType{oidc.ResponseTypeCode},
		grantTypes:                     []oidc.GrantType{oidc.GrantTypeCode, oidc.GrantTypeRefreshToken},
		accessTokenType:                op.AccessTokenTypeBearer,
		devMode:                        false,
		idTokenUserinfoClaimsAssertion: false,
		clockSkew:                      0,
	}
}

// WebClient will create a client of type web, which will always use Basic Auth and allow the use of refresh tokens
// user-defined redirectURIs may include:
// - http://localhost with port specification (e.g. http://localhost:9999/auth/callback)
// (the example will be used as default, if none is provided)
func WebClient(id, secret string, redirectURIs ...string) *Client {
	if len(redirectURIs) == 0 {
		redirectURIs = []string{
			"http://localhost:9999/auth/callback",
		}
	}
	return &Client{
		id:                             uuid.MustParse(id),
		secret:                         secret,
		redirectURIs:                   redirectURIs,
		applicationType:                op.ApplicationTypeWeb,
		authMethod:                     oidc.AuthMethodBasic,
		loginURL:                       defaultLoginURL,
		responseTypes:                  []oidc.ResponseType{oidc.ResponseTypeCode},
		grantTypes:                     []oidc.GrantType{oidc.GrantTypeCode, oidc.GrantTypeRefreshToken},
		accessTokenType:                op.AccessTokenTypeBearer,
		devMode:                        false,
		idTokenUserinfoClaimsAssertion: false,
		clockSkew:                      0,
	}
}

// DeviceClient creates a device client with Basic authentication.
func DeviceClient(id, secret string) *Client {
	return &Client{
		id:                             uuid.MustParse(id),
		secret:                         secret,
		redirectURIs:                   nil,
		applicationType:                op.ApplicationTypeWeb,
		authMethod:                     oidc.AuthMethodBasic,
		loginURL:                       defaultLoginURL,
		responseTypes:                  []oidc.ResponseType{oidc.ResponseTypeCode},
		grantTypes:                     []oidc.GrantType{oidc.GrantTypeDeviceCode},
		accessTokenType:                op.AccessTokenTypeBearer,
		devMode:                        false,
		idTokenUserinfoClaimsAssertion: false,
		clockSkew:                      0,
	}
}

type hasRedirectGlobs struct {
	*Client
}

// RedirectURIGlobs provide wildcarding for additional valid redirects
func (c hasRedirectGlobs) RedirectURIGlobs() []string {
	return c.redirectURIGlobs
}

// PostLogoutRedirectURIGlobs provide extra wildcarding for additional valid redirects
func (c hasRedirectGlobs) PostLogoutRedirectURIGlobs() []string {
	return c.postLogoutRedirectURIGlobs
}

// RedirectGlobsClient wraps the client in a op.HasRedirectGlobs
// only if DevMode is enabled.
func RedirectGlobsClient(client *Client) op.Client {
	if client.devMode {
		return hasRedirectGlobs{client}
	}
	return client
}
