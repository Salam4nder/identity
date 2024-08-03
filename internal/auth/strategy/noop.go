package strategy

import (
	"context"

	"github.com/Salam4nder/user/internal/auth"
)

var _ auth.Strategy = (*Authenticator)(nil)

func New() *Authenticator { return &Authenticator{} }

type (
	Authenticator struct{}
	Input         struct{}
	Output        struct{}
)

func (x Authenticator) Authenticate(ctx context.Context) error {
	return nil
}
