package db

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/zltl/xoidc/gen/xoidc/public/model"
	"github.com/zltl/xoidc/gen/xoidc/public/table"
	"github.com/zltl/xoidc/pkg/m"
)

func (s *Store) QueryAuthRequestByID(ctx context.Context, id string) (*m.AuthRequest, error) {
	var res model.AuthRequest
	tb := table.AuthRequest

	stmt := tb.SELECT(
		tb.AllColumns,
	).WHERE(
		tb.ID.EQ(UUID(uuid.MustParse(id))),
	).ORDER_BY(
		tb.CreationDate.DESC(),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, s.db, &res)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	a := &m.AuthRequest{
		ID:           res.ID.String(),
		CreationDate: res.CreationDate,
		UserID:       res.UserID,
		IsDone:       res.Done,
		AuthTime:     res.AuthTime,
	}

	err = a.SetContent(res.Content)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return a, nil
}

func (s *Store) StoreAuthRequest(ctx context.Context, a *m.AuthRequest) (string, error) {
	stmt := `
INSERT INTO auth_request (
    id,
    creation_date,
    user_id,
    done,
    auth_time,
    content
) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING id
`
	var uuid string
	err := s.db.QueryRowContext(
		ctx,
		stmt,
		a.CreationDate,
		a.UserID,
		a.IsDone,
		a.AuthTime,
		a.Content(),
	).Scan(&uuid)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return uuid, nil
}

func (s *Store) UpdateAuthRequest(ctx context.Context, a *m.AuthRequest) error {
	stmt := `
UPDATE auth_request
SET user_id=$1,
	done=$2,
	auth_time=$3
WHERE
    id=$4
`
	_, err := s.db.ExecContext(
		ctx,
		stmt,
		a.UserID,
		a.Done(),
		a.AuthTime,
		uuid.MustParse(a.ID),
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *Store) DeleteAuthRequest(ctx context.Context, id string) error {
	tb := table.AuthRequest
	stmt := tb.DELETE().WHERE(
		tb.ID.EQ(UUID(uuid.MustParse(id))),
	)

	_, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
