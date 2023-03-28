package main

import (
	"context"
	"time"

	"github.com/Salam4nder/inventory/pkg/logger"
	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/grpc"
	"github.com/Salam4nder/user/internal/storage"
	"github.com/Salam4nder/user/pkg/mongo"
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

	mongoDB, err := mongo.New(ctx, cfg.Mongo)
	panicOnErr(err)
	defer mongoDB.Close(ctx)

	userStorage := storage.NewUserStorage(
		mongoDB.GetCollection())

	logger, err := logger.New("")

	service := grpc.NewUserService(userStorage, logger)

	server := grpc.NewServer(service, &cfg.Server, logger)

	err = server.Serve()
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
