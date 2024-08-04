package auth

import (
	"context"

	"github.com/Salam4nder/identity/proto/gen"
)

type Strategy interface {
	ConfiguredStrategy() gen.Strategy
	Revoke(context.Context) error
	Register(context.Context) error
	Authenticate(context.Context) error
}
