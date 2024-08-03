package strategy

import (
	"context"
)

type (
	Credentials struct{}
)

// var _ auth.Strategy = (*Credentials)(nil)

func (x Credentials) Authenticate(ctx context.Context) error
