package auth

import "context"

type Authenticator interface {
	Authenticate(ctx context.Context, in Input) (Output, error)
}

type Input interface {
	Valid() error
}

type Output interface {
}
