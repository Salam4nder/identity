package server

import (
	"context"
	"fmt"

	"github.com/Salam4nder/identity/internal/auth/strategy/credentials"
	"github.com/Salam4nder/identity/internal/auth/strategy/personalnumber"
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

	registerResponse := new(gen.RegisterResponse)
	var (
		err        error
		requestCtx context.Context
	)
	switch strat {
	case gen.Strategy_TypeCredentials:
		ctx = credentials.NewContext(ctx, &credentials.Input{
			Email:    req.GetCredentials().GetEmail(),
			Password: req.GetCredentials().GetPassword(),
		})

		requestCtx, err = x.strategies[strat].Register(ctx)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}

		creds, err := credentials.FromContext(requestCtx)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}

		registerResponse.Data = &gen.RegisterResponse_Credentials{Credentials: &gen.CredentialsOutput{Email: creds.Email}}
	case gen.Strategy_TypePersonalNumber:
		requestCtx, err = x.strategies[strat].Register(ctx)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}

		n, err := personalnumber.FromContext(requestCtx)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}

		registerResponse.Data = &gen.RegisterResponse_Number{Number: &gen.PersonalNumber{Numbers: n}}

	default:
		return nil, internalServerError(ctx, fmt.Errorf("unsupported strategy %s", req.GetStrategy().String()))
	}

	metrics.UsersActive.Inc()
	metrics.UsersRegistered.Inc()

	return registerResponse, nil
}
