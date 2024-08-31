package auth

import (
	"context"
	"errors"

	"github.com/Salam4nder/identity/internal/auth/strategy"
	"github.com/Salam4nder/identity/proto/gen"
)

type Strategy interface {
	// Renew will trade a valid refresh token for a new access token.
	Renew(context.Context) error
	// Revoke will purge all active tokens in the configured hot-storage.
	Revoke(context.Context) error
	// Register an entry with the configured strategy.
	Register(context.Context) error
	// Authenticate the user with the configured strategy.
	Authenticate(context.Context) error
}

const (
	StrategyCredentials    = "credentials"
	StrategyPersonalNumber = "personal_number"
)

var (
	_ Strategy = (*strategy.Credentials)(nil)
	_ Strategy = (*strategy.PersonalNumber)(nil)
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
