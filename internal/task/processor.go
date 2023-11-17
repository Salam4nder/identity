package task

import (
	"context"

	"github.com/Salam4nder/user/internal/db"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
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

// NewRedisTaskProcessor creates a new redis processor.
func NewRedisTaskProcessor(db db.Storage, redisOpt asynq.RedisClientOpt) Processor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  1,
		},
	})

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
