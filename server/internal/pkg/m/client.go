package m

import (
	"fmt"

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
	ID                             string   `json:"id"`
	Secret                         string   `json:"secret"`
	RedirectURIs                   []string `json:"redirect_uris"`
	ApplicationType                int      `json:"application_type"`
	AuthMethod                     string   `json:"auth_method"`
	ResponseTypes                  []string `json:"response_types"`
	GrantTypes                     []string `json:"grant_types"`
	AccessTokenType                int      `json:"access_token_type"`
	DevMode                        bool     `json:"dev_mode"`
	IDTokenUserinfoClaimsAssertion bool     `json:"id_token_userinfo_claims_assertion"`
	ClockSkew                      string   `json:"clock_skew"`
	PostLogoutRedirectURIGlobs     []string `json:"post_logout_redirect_uri_globs"`
	RedirectURIGlobs               []string `json:"redirect_uri_globs"`
	UserNamespaceID                string   `json:"user_namespace_id"`
	Name                           string   `json:"name"`
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
