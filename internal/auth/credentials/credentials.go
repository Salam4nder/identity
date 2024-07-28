package credentials

import "context"

type (
	Authenticator struct{}
	Input         struct{}
	Output        struct{}
)

func (x Authenticator) Authenticate(ctx context.Context, in Input) (Output, error)

func (x Input) Valid() error { return nil }
