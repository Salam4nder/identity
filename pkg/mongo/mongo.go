package mongo

import (
	"context"
	"time"

	"github.com/Salam4nder/user/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout = 10 * time.Second
	maxConnIdleTime   = 3 * time.Minute
	minPoolSize       = 10
	maxPoolSize       = 100
)

// DB implements the mongo.Database interface.
type DB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// New creates a new MongoDB instance.
func New(ctx context.Context, cfg config.MongoDB) (*DB, error) {
	opts := options.Client().ApplyURI(cfg.URI()).
		SetAuth(
			options.Credential{
				Username: cfg.Username,
				Password: cfg.Password,
			},
		).SetConnectTimeout(connectionTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database(cfg.Name).Collection(cfg.Collection)

	return &DB{
		client:     client,
		collection: collection,
	}, nil
}

// GetCollection returns the collection of the database.
func (db *DB) GetCollection() *mongo.Collection {
	return db.collection
}

// Close closes the database connection.
func (db *DB) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

// Ping pings the database to check the connection.
func (db *DB) Ping(ctx context.Context) error {
	return db.client.Ping(ctx, nil)
}
