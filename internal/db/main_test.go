//go:build testdb
// +build testdb

package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Salam4nder/user/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stimtech/go-migration"
	"go.uber.org/zap"
)

const (
	connTimeout     = 15 * time.Second
	migrationFolder = "migrations"
)

var (
	ctx             = context.Background()
	TestSQLConnPool *SQL
)

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Postgres{
		Host:     "localhost",
		Port:     "54321",
		Name:     "unit-test-user-db",
		User:     "test",
		Password: "unit-test-pw",
	}

	db, err := sql.Open(cfg.Driver(), cfg.Addr())
	if err != nil {
		log.Error().Err(err).
			Msg("db main_test: failed to connect to db, try running make test-db from project root")

		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), connTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Error().Err(err).
			Msg("db main_test: failed to ping db, try running make test-db from project root")

		os.Exit(1)
	}

	if err := migration.New(db, zap.NewNop()).WithFolder(migrationFolder).Migrate(); err != nil {
		log.Error().Err(err).
			Msg("db main_test: failed to migrate db, try running make test-db from project root")

		os.Exit(1)
	}

	TestSQLConnPool = &SQL{db: db}

	log.Info().Msg("db main_test: successfully connected to db")

	os.Exit(m.Run())
}

// NewTestSQLConnPool returns a new SQL connection pool for testing.
// [tablename] is the name of the table to truncate after each test.
func NewTestSQLConnPool(tablename string) (*SQL, func()) {
	return TestSQLConnPool, func() {
		_, err := TestSQLConnPool.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", tablename))
		if err != nil {
			log.Error().Err(err).
				Msg("db main_test: failed to truncate db")
		}
	}
}
