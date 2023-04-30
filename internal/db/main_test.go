package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Salam4nder/user/internal/config"
)

func TestMain(m *testing.M) {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)

		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = NewSQLDatabase(ctx, cfg.PSQL)
	if err != nil {
		log.Println(err)

		os.Exit(1)
	}

	os.Exit(m.Run())
}
