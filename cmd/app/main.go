package main

import (
	"context"
	"time"

	"github.com/Salam4nder/inventory/pkg/logger"
	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/grpc"
	"github.com/Salam4nder/user/internal/storage"
	"github.com/Salam4nder/user/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	initDBIndexes(ctx, mongoDB.GetCollection())

	userStorage := storage.NewUserStorage(
		mongoDB.GetCollection())

	logger, err := logger.New("")

	service, err := grpc.NewUserService(
		userStorage, logger, cfg.Service)
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

// initDBIndexes creates indexes for the collections.
func initDBIndexes(ctx context.Context, colls ...*mongoDB.Collection) error {
	indexModel := mongoDB.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	for _, coll := range colls {
		if coll.Name() == "users" {
			_, err := coll.Indexes().CreateOne(ctx, indexModel)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
