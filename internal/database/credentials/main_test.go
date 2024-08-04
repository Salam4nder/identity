//go:build testdb
// +build testdb

package credentials_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/database/credentials"
)

var testConn *sql.DB

func Conn() (*sql.DB, func()) {
	return testConn, func() {
		_, err := testConn.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", credentials.Tablename))
		if err != nil {
			slog.Error(fmt.Sprintf("truncating table %s", credentials.Tablename), "err", err)
		}
	}
}

func TestMain(m *testing.M) {
	cfg := config.PSQLTestConfig()

	db, err := sql.Open(cfg.Driver(), cfg.Addr())
	if err != nil {
		slog.Error("database: opening sql", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("database: pinging", "err", err)
		os.Exit(1)
	}

	testConn = db
	os.Exit(m.Run())
}
