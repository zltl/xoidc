package storage

import (
	"context"

	"encoding/json"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zltl/xoidc/server/gen/xoidc/public/model"
	"github.com/zltl/xoidc/server/gen/xoidc/public/table"
)

func (s *Storage) TXGetAuthRequestByUUID(ctx context.Context, tx qrm.DB, id uuid.UUID) (*AuthRequest, error) {
	var res model.AuthRequest
	tb := table.AuthRequest

	stmt := tb.SELECT(
		tb.AllColumns,
	).WHERE(
		tb.ID.EQ(UUID(id)),
	).ORDER_BY(
		tb.CreationDate.DESC(),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, tx, &res)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	a := &AuthRequest{
		ID:           res.ID,
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

func (s *Storage) TXStoreAuthRequest(ctx context.Context, tx qrm.DB, a *AuthRequest) (uuid.UUID, error) {
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
	var uid uuid.UUID
	rows, err := tx.QueryContext(
		ctx,
		stmt,
		a.CreationDate,
		a.UserID,
		a.IsDone,
		a.AuthTime,
		a.Content(),
	)
	if err != nil {
		logrus.Error(err)
		return uuid.UUID{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&uid)
		if err != nil {
			logrus.Error(err)
			return uuid.UUID{}, err
		}
		break
	}

	return uid, nil
}

func (s *Storage) TXUpdateAuthRequest(ctx context.Context, tx qrm.DB, a *AuthRequest) error {
	stmt := `
UPDATE auth_request
SET user_id=$1,
    done=$2,
    auth_time=$3
WHERE
    id=$4
`
	_, err := tx.ExecContext(
		ctx,
		stmt,
		a.UserID,
		a.Done(),
		a.AuthTime,
		a.ID,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *Storage) GetAuthRequestByUUID(ctx context.Context, id uuid.UUID) (*AuthRequest, error) {
	return s.TXGetAuthRequestByUUID(ctx, s.db, id)
}

func (s *Storage) StoreAuthRequest(ctx context.Context, a *AuthRequest) (uuid.UUID, error) {
	return s.TXStoreAuthRequest(ctx, s.db, a)
}

func (s *Storage) UpdateAuthRequest(ctx context.Context, a *AuthRequest) error {
	return s.TXUpdateAuthRequest(ctx, s.db, a)
}

func (s *Storage) DeleteAuthRequestByUUID(ctx context.Context, id uuid.UUID) error {
	tb := table.AuthRequest
	stmt := tb.DELETE().WHERE(
		tb.ID.EQ(UUID(id)),
	)

	_, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

const (
	// CustomScope is an example for how to use custom scopes in this library
	//(in this scenario, when requested, it will return a custom claim)
	CustomScope = "custom_scope"

	// CustomClaim is an example for how to return custom claims with this library
	CustomClaim = "custom_claim"

	// CustomScopeImpersonatePrefix is an example scope prefix for passing user id to impersonate using token exchage
	CustomScopeImpersonatePrefix = "custom_scope:impersonate:"
)

type AuthRequest struct {
	ID           uuid.UUID
	CreationDate time.Time
	AuthReq      oidc.AuthRequest
	UserID       uuid.UUID
	IsDone       bool
	AuthTime     time.Time
}

func (a *AuthRequest) GetID() string {
	return a.ID.String()
}

func (a *AuthRequest) GetACR() string {
	return "" // we won't handle acr in this example
}

func (a *AuthRequest) GetAMR() []string {
	// this example only uses password for authentication
	if a.IsDone {
		return []string{"pwd"}
	}
	return nil
}

func (a *AuthRequest) GetAudience() []string {
	return []string{a.AuthReq.ClientID} // always just use the client_id as audience
}

func (a *AuthRequest) GetAuthTime() time.Time {
	return a.AuthTime
}

func (a *AuthRequest) GetClientID() string {
	return a.AuthReq.ClientID
}

func (a *AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(&OIDCCodeChallenge{
		Challenge: a.AuthReq.CodeChallenge,
		Method:    string(a.AuthReq.CodeChallengeMethod),
	})
}

func (a *AuthRequest) GetNonce() string {
	return a.AuthReq.Nonce
}

func (a *AuthRequest) GetRedirectURI() string {
	return a.AuthReq.RedirectURI
}

func (a *AuthRequest) GetResponseType() oidc.ResponseType {
	return a.AuthReq.ResponseType
}

func (a *AuthRequest) GetResponseMode() oidc.ResponseMode {
	return "" // we won't handle response mode in this example
}

func (a *AuthRequest) GetScopes() []string {
	return a.AuthReq.Scopes
}

func (a *AuthRequest) GetState() string {
	return a.AuthReq.State
}

func (a *AuthRequest) GetSubject() string {
	return a.UserID.String()
}

func (a *AuthRequest) Done() bool {
	return a.IsDone
}

func (a *AuthRequest) Content() string {
	// json a.AuthReq
	js, err := json.Marshal(a.AuthReq)
	if err != nil {
		logrus.Fatal(err)
	}
	return string(js)
}

func (a *AuthRequest) SetContent(ct string) error {
	var arq oidc.AuthRequest
	err := json.Unmarshal([]byte(ct), &arq)
	if err != nil {
		return err
	}
	a.AuthReq = arq
	return nil
}
