package storage

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/server/gen/xoidc/public/table"
)

type Token struct {
	ID             uuid.UUID
	ApplicationID  uuid.UUID
	Subject        uuid.UUID
	RefreshTokenID uuid.UUID
	Audience       []string
	Expiration     time.Time
	Scopes         []string
}

type RefreshToken struct {
	ID            uuid.UUID
	Token         string
	AuthTime      time.Time
	AMR           []string
	Audience      []string
	UserID        uuid.UUID
	ApplicationID uuid.UUID
	Expiration    time.Time
	Scopes        []string
}

func (s *Storage) SaveToken(ctx context.Context, token *Token) error {
	tb := table.Token
	stmt := tb.INSERT(
		tb.ID,
		tb.ApplicationID,
		tb.Subject,
		tb.RefreshTokenID,
		tb.Audience,
		tb.Expiration,
		tb.Scopes,
	).VALUES(
		token.ID,
		token.ApplicationID,
		token.Subject,
		token.RefreshTokenID,
		pq.Array(token.Audience),
		token.Expiration,
		pq.Array(token.Scopes),
	)
	cmd, args := stmt.Sql()
	_, err := s.db.ExecContext(ctx, cmd, args...)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *Storage) StoreRefreshToken(ctx context.Context, reftok *RefreshToken) error {
	tb := table.RefreshToken
	stmt := tb.INSERT(
		tb.ID,
		tb.Token,
		tb.AuthTime,
		tb.Amr,
		tb.Audience,
		tb.UserID,
		tb.ApplicationID,
		tb.Expiration,
		tb.Scopes,
	).VALUES(
		reftok.ID,
		reftok.Token,
		reftok.AuthTime,
		pq.Array(reftok.AMR),
		pq.Array(reftok.Audience),
		reftok.UserID,
		reftok.ApplicationID,
		reftok.Expiration,
		pq.Array(reftok.Scopes),
	)
	cmd, args := stmt.Sql()
	_, err := s.db.ExecContext(ctx, cmd, args...)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *Storage) QueryToken(ctx context.Context, id uuid.UUID) (Token, error) {
	cmd := `
		SELECT
			id,
			application_id,
			subject,
			refresh_token_id,
			audience,
			expiration,
			scopes
		FROM token
		WHERE id = $1
	`
	var token Token
	err := s.db.QueryRowContext(ctx, cmd, id).Scan(
		&token.ID,
		&token.ApplicationID,
		&token.Subject,
		&token.RefreshTokenID,
		pq.Array(&token.Audience),
		&token.Expiration,
		pq.Array(&token.Scopes),
	)
	if err != nil {
		logrus.Error(err)
		return Token{}, err
	}
	return token, nil
}

func (s *Storage) DeleteTokenByID(ctx context.Context, id uuid.UUID) error {
	tb := table.Token
	stmt := tb.DELETE().WHERE(
		tb.ID.EQ(postgres.UUID(id)),
	)
	cmd, args := stmt.Sql()
	_, err := s.db.ExecContext(ctx, cmd, args...)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *Storage) DeleteTokenByRefreshTokenID(ctx context.Context, id uuid.UUID) error {
	tb := table.Token
	stmt := tb.DELETE().WHERE(
		tb.RefreshTokenID.EQ(postgres.UUID(id)),
	)
	cmd, args := stmt.Sql()
	_, err := s.db.ExecContext(ctx, cmd, args...)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *Storage) DeleteRefreshTokenByApplicationAndSubject(
	ctx context.Context,
	tx qrm.DB,
	applicationID, subject uuid.UUID) error {
	cmd := `
	delete from refresh_token
	WHERE refresh_token_id in (
	select refresh_token_id from token
	where application_id=$1
	and subject=$2)
	`
	_, err := tx.ExecContext(ctx, cmd, applicationID, subject)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *Storage) DeleteTokenByApplicationAndSubject(
	ctx context.Context,
	tx qrm.DB,
	applicationID,
	subject uuid.UUID) error {
	cmd := `
		DELETE FROM token
		WHERE application_id=$1
		AND subject=$2
		`
	_, err := tx.ExecContext(ctx, cmd, applicationID, subject)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *Storage) QueryRefreshToken(ctx context.Context, id uuid.UUID) (RefreshToken, error) {
	cmd := `
		SELECT
			id,
			token,
			auth_time,
			amr,
			audience,
			user_id,
			application_id,
			expiration,
			scopes
		FROM refresh_token
		WHERE id = $1
	`
	var token RefreshToken
	err := s.db.QueryRowContext(ctx, cmd, id).Scan(
		&token.ID,
		&token.Token,
		&token.AuthTime,
		pq.Array(&token.AMR),
		pq.Array(&token.Audience),
		&token.UserID,
		&token.ApplicationID,
		&token.Expiration,
		pq.Array(&token.Scopes),
	)
	if err != nil {
		logrus.Error(err)
		return RefreshToken{}, err
	}
	return token, nil
}

func (s *Storage) DeleteRefreshTokenByID(ctx context.Context, id uuid.UUID) error {
	cmd := `
		DELETE FROM refresh_token
		WHERE id=$1
	`
	_, err := s.db.ExecContext(ctx, cmd, id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
