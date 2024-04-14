// Package db provides a wrapper around sql.DB which provides a transactional context.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

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
	PingContext(
		ctx context.Context,
		timeout time.Duration,
		interrupt chan os.Signal,
	) error

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

func New(db *sql.DB) *SQL {
	return &SQL{db: db}
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
func (x *SQL) PingContext(ctx context.Context, timeout time.Duration, interrupt chan os.Signal) error {
	log.Info().Msg("entered ping context...")
	for {
		select {
		case <-time.After(2 * time.Second):
			err := x.db.PingContext(ctx)
			if err == nil {
				return nil
			}
			log.Error().Err(err).Msgf("db: pinging database, retrying with timeout of %s...", timeout)
		case <-ctx.Done():
			return fmt.Errorf("db: pinging database: %w", ctx.Err())
		case <-interrupt:
			log.Info().Msg("db: pinging database interrupted...")
			return nil
		case <-time.After(timeout):
			return fmt.Errorf("db: pinging database timed out after %s", timeout)
		}
	}
}
