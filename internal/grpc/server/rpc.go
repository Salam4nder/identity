package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/Salam4nder/identity/internal/auth/strategy/credentials"
	"github.com/Salam4nder/identity/internal/auth/strategy/personalnumber"
	"github.com/Salam4nder/identity/internal/observability/metrics"
	"github.com/Salam4nder/identity/proto/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/emptypb"
)

var tracer = otel.Tracer("server")

func (x *Identity) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError()
	}

	strategy := req.GetStrategy()
	span.SetAttributes(attribute.String("strategy", strategy.String()))

	registerResponse := new(gen.RegisterResponse)
	var (
		err        error
		requestCtx context.Context
	)
	switch strategy {
	case gen.Strategy_TypeCredentials:
		ctx = credentials.NewContext(ctx, &credentials.Input{
			Email:    req.GetCredentials().GetEmail(),
			Password: req.GetCredentials().GetPassword(),
		})

		requestCtx, err = x.strategies[strategy].Register(ctx)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}

		creds, err := credentials.FromContext(requestCtx)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}

		registerResponse.Data = &gen.RegisterResponse_Credentials{Credentials: &gen.CredentialsOutput{Email: creds.Email}}
	case gen.Strategy_TypePersonalNumber:
		requestCtx, err = x.strategies[strategy].Register(ctx)
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

// Verify a user that registered using the credentials strategy.
func (x *Identity) Verify(ctx context.Context, req *gen.VerifyRequest) (*emptypb.Empty, error) {
	ctx, span := tracer.Start(ctx, "Verify")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError()
	}

	c, ok := x.strategies[gen.Strategy_TypeCredentials]
	if !ok {
		return nil, internalServerError(ctx, errors.New("rpc: getting strategy"))
	}
	switch credStrat := c.(type) {
	case *credentials.Strategy:
		if err := credStrat.Verify(ctx, req.GetToken()); err != nil {
			if errors.Is(err, credentials.ErrTokenDoesNotExist) {
				return nil, unauthenticatedError(ctx, err, "incorrect token")
			}
			return nil, internalServerError(ctx, err)
		}
	default:
		return nil, internalServerError(ctx, errors.New("rpc: strategy is not credentials"))
	}

	return &emptypb.Empty{}, nil
}
