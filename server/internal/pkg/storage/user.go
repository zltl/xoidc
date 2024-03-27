package storage

import (
	"context"
	"crypto/rsa"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/server/gen/xoidc/public/table"
	"golang.org/x/text/language"
)

func (s *Storage) QueryPassword(ctx context.Context, name string, namespace uuid.UUID) (string, error) {
	tb := table.User
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
	logrus.Debug(args)

	return pass.String, err
}

func (s *Storage) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
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
	u := &User{}

	var locale string
	err := s.db.QueryRowContext(ctx, cmd, args...).Scan(
		&u.Username,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.EmailVerified,
		&u.Phone,
		&u.PhoneVerified,
		&locale,
	)
	if err != nil {
		return nil, err
	}
	u.ID = id

	return u, nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, name string, clientID uuid.UUID) (*User, error) {
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
	u := &User{}
	logrus.Debugf("args=%+v", args)
	var locale string
	err := s.db.QueryRowContext(ctx, cmd, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.EmailVerified,
		&u.Phone,
		&u.PhoneVerified,
		&locale,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

type User struct {
	ID                uuid.UUID
	Username          string
	Password          string
	FirstName         string
	LastName          string
	Email             string
	EmailVerified     bool
	Phone             string
	PhoneVerified     bool
	PreferredLanguage language.Tag
	IsAdmin           bool
}

type Service struct {
	keys map[string]*rsa.PublicKey
}
