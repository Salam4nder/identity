package server

import (
	"context"
	"fmt"

	"github.com/Salam4nder/identity/internal/auth/strategy"
	"github.com/Salam4nder/identity/internal/observability/metrics"
	"github.com/Salam4nder/identity/proto/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("server")

func (x *Identity) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError()
	}

	strat := req.GetStrategy()
	span.SetAttributes(attribute.String("strategy", strat.String()))

	switch strat {
	case gen.Strategy_TypeCredentials:
		ctx = strategy.NewContext(ctx, strategy.Input{
			Email:    req.GetCredentials().GetEmail(),
			Password: req.GetCredentials().GetPassword(),
		})
	default:
		return nil, internalServerError(ctx, fmt.Errorf("unsupported strategy %s", req.GetStrategy().String()))
	}

	if err := x.strategies[strat].Register(ctx); err != nil {
		return nil, internalServerError(ctx, err)
	}

	metrics.UsersActive.Inc()
	metrics.UsersRegistered.Inc()

	return &gen.RegisterResponse{}, nil
}
