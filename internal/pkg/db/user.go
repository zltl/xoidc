package db

import (
	"context"
	"database/sql"

	"github.com/bwmarrin/snowflake"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/gen/xoidc/public/model"
	"github.com/zltl/xoidc/gen/xoidc/public/table"
)

func (s *Store) QueryPassword(ctx context.Context, name string, namespace int64) (string, error) {
	tb := table.User
	_ = namespace // TOPDO: namespace
	stmt := tb.SELECT(
		tb.Password,
	).WHERE(
		tb.Username.EQ(String(name)).AND(
			tb.Namespace.EQ(Int64(namespace)),
		),
	)
	cmd, args := stmt.Sql()

	var pass sql.NullString
	err := s.db.QueryRowContext(ctx, cmd, args...).Scan(&pass)
	log.Debug(args)

	return pass.String, err
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	sid, err := snowflake.ParseBase64(id)
	if err != nil {
		return nil, err
	}

	tb := table.User
	stmt := tb.SELECT(
		tb.Username,
		tb.Password,
		tb.GivenName,
		tb.FamilyName,
		tb.Email,
		tb.EmailVerified,
		tb.PhoneNumber,
		tb.PhoneNumberVerified,
		tb.Locale,
	).WHERE(
		tb.ID.EQ(Int64(sid.Int64())),
	)
	cmd, args := stmt.Sql()
	u := &model.User{
		ID: sid.Int64(),
	}

	err = s.db.QueryRowContext(ctx, cmd, args...).Scan(
		&u.Username,
		&u.Password,
		&u.GivenName,
		&u.FamilyName,
		&u.Email,
		&u.EmailVerified,
		&u.PhoneNumber,
		&u.PhoneNumberVerified,
		&u.Locale,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, name string) (*model.User, error) {
	tb := table.User
	stmt := tb.SELECT(
		tb.ID,
		tb.Username,
		tb.Password,
		tb.GivenName,
		tb.FamilyName,
		tb.Email,
		tb.EmailVerified,
		tb.PhoneNumber,
		tb.PhoneNumberVerified,
		tb.Locale,
	).WHERE(
		tb.Username.EQ(String(name)),
	)
	cmd, args := stmt.Sql()
	u := &model.User{}
	logrus.Debugf("args=%+v", args)
	err := s.db.QueryRowContext(ctx, cmd, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.GivenName,
		&u.FamilyName,
		&u.Email,
		&u.EmailVerified,
		&u.PhoneNumber,
		&u.PhoneNumberVerified,
		&u.Locale,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
