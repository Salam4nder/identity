package main

import (
	"context"
	"os"
	"time"

	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stimtech/go-migration"
	"go.uber.org/zap"
)

const (
	timeout = 20 * time.Second
)

func main() {
	cfg, err := config.New()
	fatalExitOnErr(err)

	var logger zerolog.Logger

	if cfg.Environment == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fatalExitOnErr(err)

		defer file.Close()

		log.Logger = log.Output(file)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	storage, err := db.NewSQLDatabase(ctx, cfg.PSQL)
	fatalExitOnErr(err)

	log.Info().Msg("successfully connected to database...")

	migration := migration.New(storage.GetDB(), zap.NewNop())

	if err := migration.Migrate(); err != nil {
		fatalExitOnErr(err)
	}
	log.Info().Msg("successfully migrated database...")

	service, err := grpc.NewUserService(storage, &logger, cfg.Service)
	fatalExitOnErr(err)

	server := grpc.NewServer(service, &cfg.Server, &logger)

	go server.ServeGRPCGateway()

	err = server.ServeGRPC()
	fatalExitOnErr(err)
}

func fatalExitOnErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("fatal exit: failed to start user service")
	}
}
