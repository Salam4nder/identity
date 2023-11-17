package task

import (
	"context"

	"github.com/hibiken/asynq"
)

// Creator is an interface for creating async tasks.
type Creator interface {
	SendVerificationEmail(
		ctx context.Context,
		payload VerificationEmailPayload,
		opts ...asynq.Option,
	) error
}

// RedisTaskCreator is a task creator that uses Redis as its implementation.
type RedisTaskCreator struct {
	client *asynq.Client
}

// NewRedisTaskCreator creates a new RedisTaskCreator.
func NewRedisTaskCreator(redisOpt asynq.RedisClientOpt) Creator {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskCreator{client: client}
}
