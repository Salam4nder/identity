package auth

import (
	"context"
	"errors"

	"github.com/Salam4nder/identity/internal/auth/strategy/credentials"
	"github.com/Salam4nder/identity/internal/auth/strategy/personalnumber"
	"github.com/Salam4nder/identity/proto/gen"
)

type Strategy interface {
	// Renew will trade a valid refresh token for a new access token.
	Renew(context.Context) error
	// Revoke will purge all active tokens in the configured hot-storage.
	Revoke(context.Context) error
	// Register an entry with the configured strategy.
	// Any output returned within context.
	Register(context.Context) (context.Context, error)
	// Authenticate the user with the configured strategy.
	Authenticate(context.Context) error
}

const (
	StrategyCredentials    = "credentials"
	StrategyPersonalNumber = "personal_number"
)

var (
	_ Strategy = (*credentials.Strategy)(nil)
	_ Strategy = (*personalnumber.Strategy)(nil)
)

func StrategyFromString(s string) (gen.Strategy, error) {
	switch s {
	case StrategyCredentials:
		return gen.Strategy_TypeCredentials, nil
	case StrategyPersonalNumber:
		return gen.Strategy_TypePersonalNumber, nil
	}
	return gen.Strategy_TypeNoStrategy, errors.New("auth: unsupported strategy")
}
