package auth

import "context"

type Strategy interface {
	Authenticate(context.Context) error
}
