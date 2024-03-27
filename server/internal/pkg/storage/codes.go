package storage

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/server/gen/xoidc/public/model"
	"github.com/zltl/xoidc/server/gen/xoidc/public/table"
)

// TODO: expire cdoe

func (s *Storage) CodeToRequestID(ctx context.Context, code string) (uuid.UUID, error) {
	tb := table.CodeRequestID
	stmt := tb.SELECT(
		tb.AllColumns,
	).WHERE(
		tb.Code.EQ(postgres.String(code)),
	).LIMIT(1)

	var mods model.CodeRequestID

	err := stmt.QueryContext(ctx, s.db, &mods)
	if err != nil {
		logrus.Error(err)
		return uuid.UUID{}, nil
	}

	return mods.RequestID, nil
}

func (s *Storage) StoreCodeRequestID(ctx context.Context, code string, requestID uuid.UUID) error {
	tb := table.CodeRequestID
	stmt := tb.INSERT(
		tb.Code,
		tb.RequestID,
	).VALUES(
		postgres.String(code),
		postgres.UUID(requestID),
	)

	cmd, args := stmt.Sql()
	_, err := s.db.ExecContext(ctx, cmd, args...)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s *Storage) DeleteCodeRequestIDByCode(ctx context.Context, code string) error {
	tb := table.CodeRequestID
	stmt := tb.DELETE().WHERE(
		tb.Code.EQ(postgres.String(code)),
	)

	cmd, args := stmt.Sql()
	_, err := s.db.ExecContext(ctx, cmd, args...)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s *Storage) DeleteCodeRequestIDByRequestID(ctx context.Context, requestID uuid.UUID) error {
	tb := table.CodeRequestID
	stmt := tb.DELETE().WHERE(
		tb.RequestID.EQ(postgres.UUID(requestID)),
	)

	cmd, args := stmt.Sql()
	_, err := s.db.ExecContext(ctx, cmd, args...)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
