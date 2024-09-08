//go:build testdb
// +build testdb

package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/Salam4nder/identity/internal/config"
	migrate "github.com/rubenv/sql-migrate"
)

const connTimeout = 15 * time.Second

func SetupTestConn(tableName string) (*sql.DB, func()) {
	cfg := config.PSQLTestConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := sql.Open(cfg.Driver(), cfg.Addr())
	if err != nil {
		slog.Error("database: opening sql", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("database: pinging", "err", err)
		os.Exit(1)
	}

	n, err := migrate.Exec(
		db,
		cfg.Driver(),
		&migrate.FileMigrationSource{Dir: "../migrations"},
		migrate.Up,
	)
	if err != nil {
		slog.Error("database: migrating", "err", err)
		os.Exit(1)
	}
	slog.Info("database: applied migrations", "amount", n)

	slog.Info("database: successfully connected to test db")
	return db, func() {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", tableName))
		if err != nil {
			slog.Error(fmt.Sprintf("truncating table %s", tableName), "err", err)
		}
	}
}
