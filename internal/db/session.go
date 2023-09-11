package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const (
	// SessionActiveByDefault is the default value for the is_active column.
	SessionActiveByDefault = true
)

// Session defines a session in the database.
type Session struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	IsActive     bool      `db:"is_active"`
	ClientIP     string    `db:"client_ip"`
	UserAgent    string    `db:"user_agent"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiresAt    time.Time `db:"expires_at"`
	RefreshToken string    `db:"refresh_token"`
}

// CreateSessionParams define the parameters to create a new session.
type CreateSessionParams struct {
	ID           uuid.UUID
	Email        string
	ClientIP     string
	UserAgent    string
	RefreshToken string
	ExpiresAt    time.Time
}

// CreateSessionTx creates a new session in the database as a transaction.
func (x *SQL) CreateSessionTx(
	ctx context.Context,
	params CreateSessionParams,
) (*Session, error) {
	var session Session

	query := `
    INSERT INTO sessions (
        id,
        email,
        is_active,
        client_ip,
        user_agent,
        created_at,
        expires_at,
        refresh_token
    ) VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    ) RETURNING
        id,
        email,
        is_active,
        client_ip,
        user_agent,
        created_at,
        expires_at,
        refresh_token
    `
	if err := x.execTx(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(
			ctx,
			query,
			params.ID,
			params.Email,
			SessionActiveByDefault,
			params.ClientIP,
			params.UserAgent,
			time.Now(),
			params.ExpiresAt,
			params.RefreshToken,
		).Scan(
			&session.ID,
			&session.Email,
			&session.IsActive,
			&session.ClientIP,
			&session.UserAgent,
			&session.CreatedAt,
			&session.ExpiresAt,
			&session.RefreshToken,
		)
	}); err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSession returns a session from the database.
func (x *SQL) GetSession(
	ctx context.Context,
	id uuid.UUID,
) (*Session, error) {
	var session Session

	query := `
    SELECT
        id,
        email,
        is_active,
        client_ip,
        user_agent,
        created_at,
        expires_at,
        refresh_token
    FROM
        sessions
    WHERE
        id = $1
    `

	if err := x.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.Email,
		&session.IsActive,
		&session.ClientIP,
		&session.UserAgent,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.RefreshToken,
	); err != nil {
		return nil, err
	}

	return &session, nil
}

// BlockSessionTx deactivates a session in the database as a transaction.
func (x *SQL) BlockSessionTx(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
    UPDATE
        sessions
    SET
        is_active = false
    WHERE
        id = $1
    `

	return x.execTx(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, query, id)
		return err
	})
}
