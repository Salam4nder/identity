// Package db provides a wrapper around sql.DB which provides a transactional context.
package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Salam4nder/user/internal/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var _ Storage = (*SQL)(nil)

// Storage has all methods to work with the database.
type Storage interface {
	// DB returns the underlying sql.DB.
	DB() *sql.DB
	// Close closes the underlying sql.DB.
	Close() error
	// PingContext pings the underlying sql.DB.
	PingContext(ctx context.Context) error

	// User repository

	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
	// ReadUser reads a user from the database.
	ReadUser(ctx context.Context, id uuid.UUID) (*User, error)
	// ReadUserByEmail reads a user from the database by email.
	ReadUserByEmail(ctx context.Context, email string) (*User, error)
	// UpdateUser updates a user in the database.
	UpdateUser(ctx context.Context, params UpdateUserParams) (*User, error)
	// DeleteUser deletes a user from the database.
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Session repository

	// CreateSession creates a new session in the database.
	CreateSession(ctx context.Context, params CreateSessionParams) (*Session, error)
	// ReadSession returns a session from the database.
	ReadSession(ctx context.Context, id uuid.UUID) (*Session, error)
	// BlockSession deactivates a session in the database.
	BlockSession(ctx context.Context, id uuid.UUID) error
}

// SQL is a wrapper around sql.DB which provides a transactional context.
type SQL struct {
	db *sql.DB
}

// DB returns the underlying sql.DB.
func (x *SQL) DB() *sql.DB {
	return x.db
}

// Close closes the underlying sql.DB.
func (x *SQL) Close() error {
	return x.db.Close()
}

// PingContext pings the underlying sql.DB.
func (x *SQL) PingContext(ctx context.Context) error {
	return x.db.PingContext(ctx)
}

// NewSQLDatabase creates a new SQLDatabase.
func NewSQLDatabase(ctx context.Context, cfg config.Postgres) (*SQL, error) {
	db, err := sql.Open(cfg.Driver(), cfg.Addr())
	if err != nil {
		return nil, fmt.Errorf("db: failed to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db: pinging database: %w", err)
	}

	log.Info().Msg("db: successfully connected to database...")

	return &SQL{db: db}, nil
}
