package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Salam4nder/user/internal/db"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// SendVerificationEmail is the task name for sending verification email.
const SendVerificationEmail = "task:send_verification_email"

// VerificationEmailPayload is the payload for sending verification email.
type VerificationEmailPayload struct {
	Email string `json:"email"`
}

// SendVerificationEmail sends a verification email task to Redis.
func (x *RedisTaskCreator) SendVerificationEmail(
	ctx context.Context,
	payload VerificationEmailPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("task: failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(SendVerificationEmail, jsonPayload, opts...)
	taskInfo, err := x.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("task: failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msg("task enqueued")

	return nil
}

// ProcessVerificationEmail processes a verification email task.
func (x *RedisTaskProcessor) ProcessVerificationEmail(ctx context.Context, task *asynq.Task) error {
	var payload VerificationEmailPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// no point in retrying since the payload is malformed.
		return fmt.Errorf("task: failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	_, err := x.db.ReadUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return fmt.Errorf("task: user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("task: failed to read user: %w", err)
	}

	// TODO: send verification email

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", payload.Email).
		Msg("verification email sent")

	return nil
}
