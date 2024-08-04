package auth

import (
	"context"

	"github.com/Salam4nder/identity/proto/gen"
)

type Strategy interface {
	// ConfiguredStrategy exposes the current configured strategy.
	ConfiguredStrategy() gen.Strategy

	// Renew will trade a valid refresh token for a new access token.
	Renew(context.Context) error
	// Revoke will purge all active tokens in the configured hot-storage.
	Revoke(context.Context) error
	// Register an entry with the configured strategy.
	Register(context.Context) error
	// Authenticate the user with the configured strategy.
	Authenticate(context.Context) error
}
