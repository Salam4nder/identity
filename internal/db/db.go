// Package db provides a wrapper around sql.DB which provides a transactional context.
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

var _ Storage = (*SQL)(nil)

// Storage has all methods to work with the database.
type Storage interface {
	// DB returns the underlying sql.DB.
	DB() *sql.DB
	// Close closes the underlying sql.DB.
	Close() error
	// PingContext pings the underlying sql.DB with n tries.
	// Returns nil if the database is reachable.
	PingContext(ctx context.Context, tries int) error

	// User repository

	// CreateUser creates a new user in the database.
	CreateUser(ctx context.Context, params CreateUserParams) error
	// ReadUser reads a user from the database.
	ReadUser(ctx context.Context, id uuid.UUID) (*User, error)
	// ReadUserByEmail reads a user from the database by email.
	ReadUserByEmail(ctx context.Context, email string) (*User, error)
	// UpdateUser updates a user in the database.
	UpdateUser(ctx context.Context, params UpdateUserParams) error
	// DeleteUser deletes a user from the database.
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Session repository

	// CreateSession creates a new session in the database.
	CreateSession(ctx context.Context, params CreateSessionParams) error
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
func (x *SQL) PingContext(ctx context.Context, maxTries int) error {
	var tries int
	for {
		if tries >= maxTries {
			return errors.New("db: pinging database, max tries reached")
		}
		select {
		case <-time.After(2 * time.Second):
			err := x.db.PingContext(ctx)
			if err == nil {
				return nil
			}
			tries++
			slog.ErrorContext(ctx, "db: pinging database", "try", tries, "maxTries", maxTries)
		case <-ctx.Done():
			return fmt.Errorf("db: pinging database: %w", ctx.Err())
		}
	}
}
