package strategy

import (
	"context"
	"log/slog"

	"github.com/Salam4nder/user/internal/auth"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("strategy")

var _ auth.Strategy = (*NoOp)(nil)

type NoOp struct{}

func NewNoOp() *NoOp { return &NoOp{} }

func (x *NoOp) Authenticate(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Authenticate")
	defer span.End()

	slog.InfoContext(ctx, "strategy: noop authentication")
	return nil
}

func (x *NoOp) Register(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	slog.InfoContext(ctx, "strategy: noop register")
	return nil
}

func (x *NoOp) Revoke(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Revoke")
	defer span.End()

	slog.InfoContext(ctx, "strategy: noop revoke")
	return nil
}
