package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// HealthCheck will ping the given db [maxTries] times.
// Returns nil when a connection is established.
func HealthCheck(ctx context.Context, db *sql.DB, maxTries int) error {
	var tries int
	for {
		if tries >= maxTries {
			return errors.New("db: pinging database, max tries reached")
		}
		select {
		case <-time.After(2 * time.Second):
			err := db.PingContext(ctx)
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
