package auth

import "context"

type Strategy interface {
	Revoke(context.Context) error
	Register(context.Context) error
	Authenticate(context.Context) error
}
