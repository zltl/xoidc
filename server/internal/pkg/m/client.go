package m

import (
	"fmt"
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"github.com/zltl/xoidc/server/internal/pkg/storage" // TODO: fix recursive import
)

type ClientListResponse struct {
	Response
	Total   int64    `json:"total"`
	Clients []Client `json:"clients"`
}

type ClientResponse struct {
	Response
	Client Client `json:"client"`
}

type Client struct {
	DID                             string   `json:"id"`
	DSecret                         string   `json:"secret"`
	DRedirectURIs                   []string `json:"redirect_uris"`
	DApplicationType                int      `json:"application_type"`
	DAuthMethod                     string   `json:"auth_method"`
	DResponseTypes                  []string `json:"response_types"`
	DGrantTypes                     []string `json:"grant_types"`
	DAccessTokenType                int      `json:"access_token_type"`
	DDevMode                        bool     `json:"dev_mode"`
	DIDTokenUserinfoClaimsAssertion bool     `json:"id_token_userinfo_claims_assertion"`
	DClockSkew                      string   `json:"clock_skew"`
	DPostLogoutRedirectURIGlobs     []string `json:"post_logout_redirect_uri_globs"`
	DRedirectURIGlobs               []string `json:"redirect_uri_globs"`
	DUserNamespaceID                string   `json:"user_namespace_id"`
	DName                           string   `json:"name"`
}

func toStrList[T any](ss []T) []string {
	r := make([]string, 0)
	for _, s := range ss {
		r = append(r, fmt.Sprint(s))
	}
	return r
}

func ClientDB2View(c *storage.Client) Client {
	return Client{
		ID:                             c.GetID(),
		Secret:                         c.GetSecret(),
		RedirectURIs:                   c.RedirectURIs(),
		ApplicationType:                int(c.ApplicationType()),
		AuthMethod:                     string(c.AuthMethod()),
		ResponseTypes:                  toStrList(c.ResponseTypes()),
		GrantTypes:                     toStrList(c.GrantTypes()),
		AccessTokenType:                int(c.AccessTokenType()),
		DevMode:                        c.DevMode(),
		IDTokenUserinfoClaimsAssertion: c.IDTokenUserinfoClaimsAssertion(),
		ClockSkew:                      c.ClockSkew().String(),
		PostLogoutRedirectURIGlobs:     c.PostLogoutRedirectURIGlobs(),
		RedirectURIGlobs:               c.RedirectURIGlobs(),
		UserNamespaceID:                c.UserNamespaceID(),
		Name:                           c.Name(),
	}
}

func (c *Client) Name() string {
	return c.DName
}

// GetID must return the client_id
func (c *Client) GetID() string {
	return c.DID
}

func (c *Client) GetSecret() string {
	return c.DSecret
}

func (c *Client) PostLogoutRedirectURIGlobs() []string {
	return c.DPostLogoutRedirectURIGlobs
}

func (c *Client) RedirectURIGlobs() []string {
	return c.DRedirectURIGlobs
}

func (c *Client) UserNamespaceID() string {
	return c.DUserNamespaceID
}

// RedirectURIs must return the registered redirect_uris for Code and Implicit Flow
func (c *Client) RedirectURIs() []string {
	return c.DRedirectURIs
}

// PostLogoutRedirectURIs must return the registered post_logout_redirect_uris for sign-outs
func (c *Client) PostLogoutRedirectURIs() []string {
	// TODO
	return []string{}
}

// ApplicationType must return the type of the client (app, native, user agent)
func (c *Client) ApplicationType() op.ApplicationType {
	return op.ApplicationType(c.DApplicationType)
}

// AuthMethod must return the authentication method (client_secret_basic, client_secret_post, none, private_key_jwt)
func (c *Client) AuthMethod() oidc.AuthMethod {
	return oidc.AuthMethod(c.DAuthMethod)
}

// ResponseTypes must return all allowed response types (code, id_token token, id_token)
// these must match with the allowed grant types
func (c *Client) ResponseTypes() []oidc.ResponseType {
	return c.DResponseTypes
}

// GrantTypes must return all allowed grant types (authorization_code, refresh_token, urn:ietf:params:oauth:grant-type:jwt-bearer)
func (c *Client) GrantTypes() []oidc.GrantType {
	return c.DGrantTypes
}

// LoginURL will be called to redirect the user (agent) to the login UI
// you could implement some logic here to redirect the users to different login UIs depending on the client
func (c *Client) LoginURL(id string) string {
	return c.loginURL(id)
}

// AccessTokenType must return the type of access token the client uses (Bearer (opaque) or JWT)
func (c *Client) AccessTokenType() op.AccessTokenType {
	return c.DAccessTokenType
}

// IDTokenLifetime must return the lifetime of the client's id_tokens
func (c *Client) IDTokenLifetime() time.Duration {
	return 1 * time.Hour
}

// DevMode enables the use of non-compliant configs such as redirect_uris (e.g. http schema for user agent client)
func (c *Client) DevMode() bool {
	return c.DDevMode
}

// RestrictAdditionalIdTokenScopes allows specifying which custom scopes shall be asserted into the id_token
func (c *Client) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}

// RestrictAdditionalAccessTokenScopes allows specifying which custom scopes shall be asserted into the JWT access_token
func (c *Client) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}

// IsScopeAllowed enables Client specific custom scopes validation
// in this example we allow the CustomScope for all clients
func (c *Client) IsScopeAllowed(scope string) bool {
	return scope == CustomScope
}

// IDTokenUserinfoClaimsAssertion allows specifying if claims of scope profile, email, phone and address are asserted into the id_token
// even if an access token if issued which violates the OIDC Core spec
// (5.4. Requesting Claims using Scope Values: https://openid.net/specs/openid-connect-core-1_0.html#ScopeClaims)
// some clients though require that e.g. email is always in the id_token when requested even if an access_token is issued
func (c *Client) IDTokenUserinfoClaimsAssertion() bool {
	return c.DIDTokenUserinfoClaimsAssertion
}

// ClockSkew enables clients to instruct the OP to apply a clock skew on the various times and expirations
// (subtract from issued_at, add to expiration, ...)
func (c *Client) ClockSkew() time.Duration {
	return c.DClockSkew
}
