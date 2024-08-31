package auth

import (
	"context"

	"github.com/Salam4nder/identity/internal/auth/strategy"
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

var (
	_ Strategy = (*strategy.Credentials)(nil)
	_ Strategy = (*strategy.PersonalNumber)(nil)
)

var supportedStrategies = map[gen.Strategy]struct{}{}

func Supported(s gen.Strategy) bool {
	if _, ok := supportedStrategies[s]; !ok {
		return false
	}
	return true
}

func strategyFromString(s string) gen.Strategy {
}

// func MountStrategies(s ...string) {
// 	for _, v := range s {
//         switch {
//             case v
//         }
// 		supportedStrategies[v] = struct{}{}
// 	}
// }
