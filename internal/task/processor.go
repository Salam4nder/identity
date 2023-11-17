package task

import (
	"context"

	"github.com/Salam4nder/user/internal/db"
	"github.com/hibiken/asynq"
)

// Processor is an interface for processing tasks.
type Processor interface {
	Process() error
	ProcessVerificationEmail(ctx context.Context, task *asynq.Task) error
}

// RedisProcessor is a processor for redis tasks.
type RedisTaskProcessor struct {
	server *asynq.Server
	db     db.Storage
}

// NewRedisProcessor creates a new redis processor.
func NewRedisProcessor(db db.Storage, redisOpt asynq.RedisClientOpt) Processor {
	server := asynq.NewServer(redisOpt, asynq.Config{})

	return &RedisTaskProcessor{
		server: server,
		db:     db,
	}
}

// Process starts the redis processor.
func (x *RedisTaskProcessor) Process() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(SendVerificationEmail, x.ProcessVerificationEmail)

	return x.server.Start(mux)
}
