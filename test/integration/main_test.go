//go:build integration

package test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Salam4nder/user/internal/config"

	"github.com/stimtech/go-migration"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)

		os.Exit(1)
	}

	db, err := sql.Open(cfg.PSQL.Driver(), cfg.PSQL.URI())
	if err != nil {
		log.Println(err)

		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Println(err)

		os.Exit(1)
	}

	_, err = db.ExecContext(ctx, Schema)
	if err != nil {
		log.Println(err)

		os.Exit(1)
	}

	migrator := migration.New(db, zap.NewNop())
	if err := migrator.Migrate(); err != nil {
		log.Println(err)

		os.Exit(1)
	}

	os.Exit(m.Run())
}
