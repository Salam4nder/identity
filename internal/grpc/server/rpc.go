package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Salam4nder/identity/internal/auth/strategy"
	"github.com/Salam4nder/identity/internal/observability/metrics"
	"github.com/Salam4nder/identity/proto/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/emptypb"
)

var tracer = otel.Tracer("server")

func (x *Identity) Register(ctx context.Context, req *gen.RegisterRequest) (*emptypb.Empty, error) {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError()
	}

	strat := req.GetStrategy()
	span.SetAttributes(attribute.String("strategy", strat.String()))

	switch strat {
	case gen.Strategy_Credentials:
		ctx = strategy.NewContext(ctx, strategy.Input{
			Email:    req.GetCredentials().GetEmail(),
			Password: req.GetCredentials().GetPassword(),
		})
	default:
		slog.ErrorContext(ctx, fmt.Sprintf("server: unsupported strategy %T,", req.GetStrategy()))
		return nil, internalServerError(ctx, fmt.Errorf("unsupported strategy %T", req.GetStrategy()))
	}

	if err := x.strategies[strat].Register(ctx); err != nil {
		return nil, internalServerError(ctx, err)
	}

	metrics.UsersActive.Inc()
	metrics.UsersRegistered.Inc()

	return &emptypb.Empty{}, nil
}
