//go:build testdb
// +build testdb

package database_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/Salam4nder/identity/internal/config"
	"github.com/stimtech/go-migration/v2"
)

const (
	connTimeout     = 15 * time.Second
	migrationFolder = "migrations"
)

// NewTestSQLConnPool returns a new SQL connection pool for testing.
// [tablename] is the name of the table to truncate for cleanup.
func TestMain(m *testing.M) {
	cfg := config.PSQLTestConfig()

	db, err := sql.Open(cfg.Driver(), cfg.Addr())
	if err != nil {
		slog.Error("database: opening sql", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), connTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("database: pinging", "err", err)
		os.Exit(1)
	}

	if err := migration.New(db, migration.Config{MigrationFolder: migrationFolder}).Migrate(); err != nil {
		slog.Error("database: migration", "err", err)
		os.Exit(1)
	}

	slog.Info("database: successfully connected to test db")
}
