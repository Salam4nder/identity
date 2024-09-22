package token

import (
	"errors"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/Salam4nder/identity/proto/gen"
)

var _ Maker = (*PasetoMaker)(nil)

const (
	identifierKey = "token_identifier"
	strategyKey   = "token_strategy"
)

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

func (x *PasetoMaker) MakeAccessToken(identifer any, strat gen.Strategy) (SafeString, error) {
	token := paseto.NewToken()
	switch strat {
	case gen.Strategy_TypeCredentials:
		s, ok := identifer.(string)
		if !ok {
			return "", fmt.Errorf("token: expected identifier to be string, got %T", identifer)
		}
		token.Set(strategyKey, gen.Strategy_TypeCredentials)
		token.Set(identifierKey, s)
	case gen.Strategy_TypePersonalNumber:
		d, ok := identifer.(uint64)
		if !ok {
			return "", fmt.Errorf("token: expected identifier to be uint64, got %T", identifer)
		}
		token.Set(strategyKey, gen.Strategy_TypePersonalNumber)
		token.Set(identifierKey, d)
	default:
		return "", errors.New("unsupported strategy")
	}
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(x.accessDur))
	return fromString(token.V4Encrypt(x.symmetricKey, nil)), nil
}

func (x *PasetoMaker) MakeRefreshToken(identifer any, strat gen.Strategy) (SafeString, error) {
	token := paseto.NewToken()
	switch strat {
	case gen.Strategy_TypeCredentials:
		s, ok := identifer.(string)
		if !ok {
			return "", fmt.Errorf("token: expected identifier to be string, got %T", identifer)
		}
		token.Set(strategyKey, gen.Strategy_TypeCredentials)
		token.Set(identifierKey, s)
	case gen.Strategy_TypePersonalNumber:
		d, ok := identifer.(uint64)
		if !ok {
			return "", fmt.Errorf("token: expected identifier to be uint64, got %T", identifer)
		}
		token.Set(strategyKey, gen.Strategy_TypePersonalNumber)
		token.Set(identifierKey, d)
	default:
		return "", errors.New("unsupported strategy")
	}
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(x.refreshDur))
	return fromString(token.V4Encrypt(x.symmetricKey, nil)), nil
}

func (x *PasetoMaker) Verify(t SafeString) error {
	_, err := x.parser.ParseV4Local(x.symmetricKey, string(t), nil)
	if err != nil {
		return fmt.Errorf("token: verifying token, %w", err)
	}
	return nil
}
