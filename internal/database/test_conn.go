//go:build testdb
// +build testdb

package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "embed"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"

	"github.com/Salam4nder/identity/internal/config"
	"github.com/Salam4nder/identity/internal/database/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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

	src, err := iofs.New(migrations.Files, ".")
	if err != nil {
		slog.Error("database: iofs", "err", err)
		os.Exit(1)
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, cfg.Addr())
	if err != nil {
		slog.Error("database: creating source", "err", err)
		os.Exit(1)
	}
	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("database: migrating", "err", err)
			os.Exit(1)
		}
	}

	return db, func() {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", tableName))
		if err != nil {
			slog.Error(fmt.Sprintf("truncating table %s", tableName), "err", err)
		}
	}
}
