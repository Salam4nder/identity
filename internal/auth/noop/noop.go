package noop

import (
	"context"

	"github.com/Salam4nder/user/internal/auth"
)

var _ auth.Authenticator = (*Authenticator)(nil)

func New() *Authenticator { return &Authenticator{} }

type (
	Authenticator struct{}
	Input         struct{}
	Output        struct{}
)

func (x Authenticator) Authenticate(ctx context.Context, in auth.Input) (auth.Output, error) {
	return Output{}, nil
}

func (x Input) Valid() error { return nil }

func (x Output) Noop()
