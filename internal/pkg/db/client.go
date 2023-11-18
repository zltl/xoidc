package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sanyokbig/pqinterval"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ID                             uuid.UUID
	Secret                         string
	RedirectURIs                   []string
	ApplicationType                int
	AuthMethod                     string
	ResponseTypes                  []string
	GrantTypes                     []string
	AccessTokenType                int
	DevMode                        bool
	IDTokenUserInfoClaimsAssertion bool
	ClockSkew                      time.Duration
	PostLogoutRedirectURIGlobs     []string
	RedirectURIGlobs               []string
	UserNamespaceID                uuid.UUID
}

func (s *Store) GetClientByID(ctx context.Context, clientID uuid.UUID) (*Client, error) {
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
			user_namespace_id
		FROM
			client
		WHERE
			id = $1
	`
	c := &Client{}
	var interval pqinterval.Interval
	err := s.db.QueryRowContext(ctx, stmt, clientID).Scan(
		&c.ID,
		&c.Secret,
		pq.Array(&c.RedirectURIs),
		&c.ApplicationType,
		&c.AuthMethod,
		pq.Array(&c.ResponseTypes),
		pq.Array(&c.GrantTypes),
		&c.AccessTokenType,
		&c.DevMode,
		&c.IDTokenUserInfoClaimsAssertion,
		&interval,
		pq.Array(&c.PostLogoutRedirectURIGlobs),
		pq.Array(&c.RedirectURIGlobs),
		&c.UserNamespaceID,
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

	c.ClockSkew = dura

	return c, err
}
