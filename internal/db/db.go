package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Salam4nder/user/internal/config"
)

// SQL is a wrapper around sql.DB which provides a transactional context.
type SQL struct {
	db *sql.DB
}

// NewSQLDatabase creates a new SQLDatabase.
func NewSQLDatabase(ctx context.Context, cfg config.Postgres) (*SQL, error) {
	db, err := sql.Open(cfg.Driver(), cfg.URI())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &SQL{db: db}, nil
}

// GetDB returns the underlying sql.DB.
func (s *SQL) GetDB() *sql.DB {
	return s.db
}

// execTx executes a function in a transaction.
func (s *SQL) execTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("rolling back transaction: %w", err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
