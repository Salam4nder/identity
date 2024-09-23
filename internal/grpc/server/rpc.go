package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Salam4nder/identity/internal/auth/strategy/credentials"
	"github.com/Salam4nder/identity/internal/auth/strategy/personalnumber"
	"github.com/Salam4nder/identity/internal/observability/metrics"
	"github.com/Salam4nder/identity/internal/token"
	"github.com/Salam4nder/identity/proto/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var tracer = otel.Tracer("server")

// Register a user with the given strategy.
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

		registerResponse.Data = &gen.RegisterResponse_Number{Number: &gen.PersonalNumber{Number: n}}

	default:
		return nil, internalServerError(ctx, fmt.Errorf("unsupported strategy %s", req.GetStrategy().String()))
	}

	metrics.UsersActive.Inc()
	metrics.UsersRegistered.Inc()

	return registerResponse, nil
}

// VerifyEmail verifies a user that registered using the credentials strategy.
func (x *Identity) VerifyEmail(ctx context.Context, req *gen.TokenRequest) (*emptypb.Empty, error) {
	ctx, span := tracer.Start(ctx, "VerifyEmail")
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
		if err := credStrat.VerifyEmail(ctx, req.GetToken()); err != nil {
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

// Authenticate a user with the given strategy.
func (x *Identity) Authenticate(ctx context.Context, req *gen.AuthenticateRequest) (*gen.AuthenticateResponse, error) {
	ctx, span := tracer.Start(ctx, "Authenticate")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError()
	}

	strategy := req.GetStrategy()
	span.SetAttributes(attribute.String("strategy", strategy.String()))

	var (
		err                       error
		accessToken, refreshToken token.SafeString
	)
	switch strategy {
	case gen.Strategy_TypeCredentials:
		ctx = credentials.NewContext(ctx, &credentials.Input{
			Email:    req.GetCredentials().GetEmail(),
			Password: req.GetCredentials().GetPassword(),
		})

		if err = x.strategies[strategy].Authenticate(ctx); err != nil {
			switch {
			case errors.Is(err, credentials.ErrUserNotFound), errors.Is(err, credentials.ErrIncorrectPassword):
				return nil, invalidArgumentError(ctx, err, err.Error())
			case errors.Is(err, credentials.ErrUserNotVerified):
				return nil, notFoundError(ctx, err, err.Error())
			default:
				return nil, internalServerError(ctx, err)
			}
		}
		accessToken, err = x.tokenMaker.MakeAccessToken(req.GetCredentials().Email, gen.Strategy_TypeCredentials)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}
		refreshToken, err = x.tokenMaker.MakeRefreshToken(req.GetCredentials().Email, gen.Strategy_TypeCredentials)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}

	case gen.Strategy_TypePersonalNumber:
		ctx = personalnumber.NewContext(ctx, req.GetNumber().GetNumber())

		if err = x.strategies[strategy].Authenticate(ctx); err != nil {
			switch {
			case errors.Is(err, personalnumber.ErrNumberNotFound):
				return nil, invalidArgumentError(ctx, err, err.Error())
			default:
				return nil, internalServerError(ctx, err)
			}
		}

		accessToken, err = x.tokenMaker.MakeAccessToken(
			req.GetNumber().GetNumber(),
			gen.Strategy_TypePersonalNumber,
		)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}
		refreshToken, err = x.tokenMaker.MakeRefreshToken(
			req.GetNumber().GetNumber(),
			gen.Strategy_TypePersonalNumber,
		)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}

	default:
		return nil, internalServerError(ctx, fmt.Errorf("unsupported strategy %s", req.GetStrategy().String()))
	}
	return &gen.AuthenticateResponse{
		AccessToken:  string(accessToken),
		RefreshToken: string(refreshToken),
	}, nil
}

// Refresh will exchange a valid refresh token for a new access token.
func (x *Identity) Refresh(ctx context.Context, req *gen.TokenRequest) (*gen.RefreshResponse, error) {
	ctx, span := tracer.Start(ctx, "Refresh")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError()
	}

	t, err := x.tokenMaker.Parse(req.GetToken())
	if err != nil {
		return nil, internalServerError(ctx, err)
	}

	exp, err := t.GetExpiration()
	if err != nil {
		return nil, internalServerError(ctx, err)
	}
	if time.Now().After(exp) {
		return nil, invalidArgumentError(ctx, errors.New("rpc: token is expired"), "token is expired")
	}

	var strat gen.Strategy
	if err = t.Get(token.PasetoStrategyKey, &strat); err != nil {
		return nil, internalServerError(ctx, err)
	}

	var accessToken token.SafeString
	switch strat {
	case gen.Strategy_TypeCredentials:
		var email string
		if err = t.Get(token.PasetoIdentifierKey, &email); err != nil {
			return nil, internalServerError(ctx, err)
		}
		accessToken, err = x.tokenMaker.MakeAccessToken(email, gen.Strategy_TypeCredentials)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}
	case gen.Strategy_TypePersonalNumber:
		var number uint64
		if err = t.Get(token.PasetoIdentifierKey, &number); err != nil {
			return nil, internalServerError(ctx, err)
		}
		accessToken, err = x.tokenMaker.MakeAccessToken(number, gen.Strategy_TypeCredentials)
		if err != nil {
			return nil, internalServerError(ctx, err)
		}
	}

	return &gen.RefreshResponse{
		Token:     string(accessToken),
		ExpiresAt: timestamppb.New(x.tokenMaker.AccessTokenExpiration()),
	}, nil
}
