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

// Storage has all methods to work with the database.
type Storage interface {
	DB() *sql.DB
	Close() error
	PingContext(ctx context.Context) error

	// User repository
	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
	ReadUser(ctx context.Context, id uuid.UUID) (*User, error)
	ReadUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, params UpdateUserParams) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Session repository
	CreateSession(ctx context.Context, params CreateSessionParams) (*Session, error)
	ReadSession(ctx context.Context, id uuid.UUID) (*Session, error)
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
func NewSQLDatabase(ctx context.Context, cfg config.Postgres) (Storage, error) {
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

// execTx executes a function in a transaction.
// Need to figure out a better implementation for this.
// func (x *SQL) execTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
// 	tx, err := x.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("db: beginning transaction: %w", err)
// 	}

// 	if err := fn(tx); err != nil {
// 		if err := tx.Rollback(); err != nil {
// 			return fmt.Errorf("db: rolling back transaction: %w", err)
// 		}
// 		return err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return fmt.Errorf("db: committing transaction: %w", err)
// 	}

// 	return nil
// }
