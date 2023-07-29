package db

import (
	"database/sql"

	"github.com/zltl/xoidc/.gen/table"
)

func (s *Store) QueryPassword(name string, namespace int) (string, error) {
	tb := table.User
	stm := tb.Select(tb.ID,
		tb.Password,
	).Where(
		tb.Name.Eq(name),
		tb.Namespace.Eq(namespace),
	)

	var pass sql.NullString
	err := stm.Query(s.db, &pass)

	return pass.String, err
}
