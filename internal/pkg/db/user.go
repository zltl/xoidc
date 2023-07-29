package db

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/zltl/xoidc/gen/xoidc/public/table"
)

func (s *Store) QueryPassword(ctx context.Context, name string, namespace int64) (string, error) {
	tb := table.User
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

	return pass.String, err
}
