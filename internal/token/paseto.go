package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

var _ Maker = (*PasetoMaker)(nil)

// PasetoMaker makes PASETO tokens.
type PasetoMaker struct {
	accessDur    time.Duration
	refreshDur   time.Duration
	symmetricKey paseto.V4SymmetricKey
	parser       *paseto.Parser
}

func BootstrapPasetoMaker(
	accessDur, refreshDur time.Duration,
	symmetricKey []byte,
) (*PasetoMaker, error) {
	k, err := paseto.V4SymmetricKeyFromBytes(symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("token: creating symmetric key, %w", err)
	}

	p := paseto.MakeParser([]paseto.Rule{
		paseto.NotExpired(),
		paseto.ValidAt(time.Now()),
	},
	)

	return &PasetoMaker{
		accessDur:    accessDur,
		refreshDur:   refreshDur,
		symmetricKey: k,
		parser:       &p,
	}, nil
}

func (x *PasetoMaker) MakeAccessToken() SafeString {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(x.accessDur))
	return fromString(token.V4Encrypt(x.symmetricKey, nil))
}

func (x *PasetoMaker) MakeRefreshToken() SafeString {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(x.refreshDur))
	return fromString(token.V4Encrypt(x.symmetricKey, nil))
}

func (x *PasetoMaker) Verify(t SafeString) error {
	_, err := x.parser.ParseV4Local(x.symmetricKey, string(t), nil)
	if err != nil {
		return fmt.Errorf("token: verifying token, %w", err)
	}
	return nil
}
