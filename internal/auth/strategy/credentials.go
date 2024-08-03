package strategy

import (
	"context"

	"github.com/Salam4nder/user/internal/auth"
)

type (
	Authenticator struct{}
	Input         struct{}
	Output        struct{}
)

var _ auth.Strategy = (*Authenticator)(nil)

func (x Authenticator) Authenticate(ctx context.Context, in Input) (Output, error)

func (x Input) Valid() error { return nil }
