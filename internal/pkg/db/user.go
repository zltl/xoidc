package db

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/gen/xoidc/public/model"
	"github.com/zltl/xoidc/gen/xoidc/public/table"
)

func (s *Store) QueryPassword(ctx context.Context, name string, namespace uuid.UUID) (string, error) {
	tb := table.User
	_ = namespace // TOPDO: namespace
	stmt := tb.SELECT(
		tb.Password,
	).WHERE(
		tb.Username.EQ(String(name)).AND(
			tb.NamespaceID.EQ(UUID(namespace)),
		),
	)
	cmd, args := stmt.Sql()

	var pass sql.NullString
	err := s.db.QueryRowContext(ctx, cmd, args...).Scan(&pass)
	log.Debug(args)

	return pass.String, err
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
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
		tb.ID.EQ(UUID(id)),
	)
	cmd, args := stmt.Sql()
	u := &model.User{}

	err := s.db.QueryRowContext(ctx, cmd, args...).Scan(
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
	u.ID = id

	return u, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, name string, clientID uuid.UUID) (*model.User, error) {
	stmt := SELECT(
		table.User.ID,
		table.User.Username,
		table.User.Password,
		table.User.GivenName,
		table.User.FamilyName,
		table.User.Email,
		table.User.EmailVerified,
		table.User.PhoneNumber,
		table.User.PhoneNumberVerified,
		table.User.Locale,
	).FROM(
		table.User,
		table.Client,
	).WHERE(
		AND(
			table.User.NamespaceID.EQ(table.Client.UserNamespaceID),
			table.User.Username.EQ(String(name)),
			table.Client.ID.EQ(UUID(clientID)),
		),
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
