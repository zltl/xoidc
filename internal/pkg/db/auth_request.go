package db

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sirupsen/logrus"

	"github.com/zltl/xoidc/gen/xoidc/public/model"
	"github.com/zltl/xoidc/gen/xoidc/public/table"
)

func (s *Store) QueryAuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error) {
	var res model.AuthRequest
	tb := table.AuthRequest

	stmt := tb.SELECT(
		tb.AllColumns,
	).WHERE(
		tb.ID.EQ(String(id)),
	).ORDER_BY(
		tb.CreationDate.DESC(),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, s.db, &res)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &res, nil
}

func (s *Store) SetAuthRequest(ctx context.Context, req *model.AuthRequest) error {
	tb := table.AuthRequest
	stmt := tb.INSERT(
		tb.AllColumns,
	).MODEL(req).ON_CONFLICT(
		tb.ID,
	).DO_UPDATE(
		SET(
			tb.CreationDate.SET(tb.EXCLUDED.CreationDate), // EXCLUDE references proposed insertion row
			tb.ApplicationID.SET(tb.EXCLUDED.ApplicationID),
			tb.CallbackURI.SET(tb.EXCLUDED.CallbackURI),
			tb.TransferState.SET(tb.EXCLUDED.TransferState),
			tb.Prompt.SET(tb.EXCLUDED.Prompt),
			tb.UILocales.SET(tb.EXCLUDED.UILocales),
			tb.LoginHint.SET(tb.EXCLUDED.LoginHint),
			tb.MaxAuthAge.SET(tb.EXCLUDED.MaxAuthAge),
			tb.UserID.SET(tb.EXCLUDED.UserID),
			tb.Scopes.SET(tb.EXCLUDED.Scopes),
			tb.ResponseType.SET(tb.EXCLUDED.ResponseType),
			tb.Nonce.SET(tb.EXCLUDED.Nonce),
			tb.OidcCodeChallange.SET(tb.EXCLUDED.OidcCodeChallange),
			tb.OidcCodeChallangeMethod.SET(tb.EXCLUDED.OidcCodeChallangeMethod),
			tb.Done.SET(tb.EXCLUDED.Done),
			tb.AuthTime.SET(tb.EXCLUDED.AuthTime),
		),
	)

	_, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s *Store) DeleteAuthRequest(ctx context.Context, id string) error {
	tb := table.AuthRequest
	stmt := tb.DELETE().WHERE(
		tb.ID.EQ(String(id)),
	)

	_, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
