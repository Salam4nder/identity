package auth

import "context"

type Authenticator interface {
	Authenticate(context.Context, Input) (Output, error)
}

type Input interface {
	Valid() error
}

type Output interface {
	Noop()
}
