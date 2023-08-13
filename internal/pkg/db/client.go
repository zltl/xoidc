package db

import (
	"context"

	"github.com/bwmarrin/snowflake"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/zltl/xoidc/gen/xoidc/public/model"
	"github.com/zltl/xoidc/gen/xoidc/public/table"
)

func (s *Store) GetClientGrantTypes(ctx context.Context, cid string) ([]string, error) {
	cidis, err := snowflake.ParseBase64(cid)
	if err != nil {
		return nil, err
	}
	cidi := cidis.Int64()

	stmt := table.ClientGrantTypes.SELECT(
		table.ClientGrantTypes.GrantType,
	).WHERE(
		table.ClientGrantTypes.ClientID.EQ(Int64(cidi)),
	)
	cmd, args := stmt.Sql()
	rows, err := s.db.QueryContext(ctx, cmd, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grantTypes []string
	for rows.Next() {
		var grantType string
		err = rows.Scan(&grantType)
		if err != nil {
			return nil, err
		}
		grantTypes = append(grantTypes, grantType)
	}
	return grantTypes, nil
}

func (s *Store) GetClientRedirectURIs(ctx context.Context, cid string) ([]string, error) {
	cidis, err := snowflake.ParseBase64(cid)
	if err != nil {
		return nil, err
	}
	cidi := cidis.Int64()

	stmt := table.ClientRedirectUris.SELECT(
		table.ClientRedirectUris.RedirectURI,
	).WHERE(
		table.ClientRedirectUris.ClientID.EQ(Int64(cidi)),
	)
	cmd, args := stmt.Sql()
	rows, err := s.db.QueryContext(ctx, cmd, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uris []string
	for rows.Next() {
		var uri string
		err = rows.Scan(&uri)
		if err != nil {
			return nil, err
		}
		uris = append(uris, uri)
	}
	return uris, nil
}

func (s *Store) GetClientResponseTypes(ctx context.Context, cid string) ([]string, error) {
	cidis, err := snowflake.ParseBase64(cid)
	if err != nil {
		return nil, err
	}
	cidi := cidis.Int64()

	stmt := table.ClientResponseTypes.SELECT(
		table.ClientResponseTypes.ResponseType,
	).WHERE(
		table.ClientResponseTypes.ClientID.EQ(Int64(cidi)),
	)
	cmd, args := stmt.Sql()
	rows, err := s.db.QueryContext(ctx, cmd, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responseTypes []string
	for rows.Next() {
		var responseType string
		err = rows.Scan(&responseType)
		if err != nil {
			return nil, err
		}
		responseTypes = append(responseTypes, responseType)
	}
	return responseTypes, nil
}

func (s *Store) GetClientByID(ctx context.Context, cid string) (*model.Client, error) {
	iid, err := snowflake.ParseBase64(cid)
	if err != nil {
		return nil, err
	}

	stmt := table.Client.SELECT(
		table.Client.ID,
		table.Client.Secret,
		table.Client.ApplicationType,
		table.Client.AuthMethod,
		table.Client.AccessTokenType,
		table.Client.DevMode,
		table.Client.IDTokenUserinfoClaimsAssertion,
		table.Client.ClockSkew,
	).WHERE(
		table.Client.ID.EQ(Int64(iid.Int64())),
	)
	cmd, args := stmt.Sql()
	c := &model.Client{}
	err = s.db.QueryRowContext(ctx, cmd, args...).Scan(
		&c.ID,
		&c.Secret,
		&c.ApplicationType,
		&c.AuthMethod,
		&c.AccessTokenType,
		&c.DevMode,
		&c.IDTokenUserinfoClaimsAssertion,
		&c.ClockSkew,
	)
	// TODO: load grant types
	// TODO: load redirect uris
	// TODO: load response types

	return c, err
}
