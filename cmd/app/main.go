package main

import (
	"context"
	"time"

	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc"

	"github.com/Salam4nder/inventory/pkg/logger"
	"github.com/stimtech/go-migration"
)

const (
	timeout = 20 * time.Second
)

func main() {
	cfg, err := config.New()
	panicOnErr(err)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	logger, err := logger.New("")

	storage, err := db.NewSQLDatabase(ctx, cfg.PSQL)
	if err != nil {
		panicOnErr(err)
	}
	logger.Info("Connected to database...")

	migration := migration.New(storage.GetDB(), logger)

	if err := migration.Migrate(); err != nil {
		panicOnErr(err)
	}

	service, err := grpc.NewUserService(storage, logger, cfg.Service)
	panicOnErr(err)

	server := grpc.NewServer(service, &cfg.Server, logger)

	go server.ServeGRPCGateway()

	err = server.ServeGRPC()
	panicOnErr(err)

}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
