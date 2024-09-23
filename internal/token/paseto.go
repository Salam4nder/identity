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
	PasetoIdentifierKey = "token_identifier"
	PasetoStrategyKey   = "token_strategy"
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
		token.Set(PasetoStrategyKey, gen.Strategy_TypeCredentials)
		token.Set(PasetoIdentifierKey, s)
	case gen.Strategy_TypePersonalNumber:
		d, ok := identifer.(uint64)
		if !ok {
			return "", fmt.Errorf("token: expected identifier to be uint64, got %T", identifer)
		}
		token.Set(PasetoStrategyKey, gen.Strategy_TypePersonalNumber)
		token.Set(PasetoIdentifierKey, d)
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
		token.Set(PasetoStrategyKey, gen.Strategy_TypeCredentials)
		token.Set(PasetoIdentifierKey, s)
	case gen.Strategy_TypePersonalNumber:
		d, ok := identifer.(uint64)
		if !ok {
			return "", fmt.Errorf("token: expected identifier to be uint64, got %T", identifer)
		}
		token.Set(PasetoStrategyKey, gen.Strategy_TypePersonalNumber)
		token.Set(PasetoIdentifierKey, d)
	default:
		return "", errors.New("unsupported strategy")
	}
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(x.refreshDur))
	return fromString(token.V4Encrypt(x.symmetricKey, nil)), nil
}

// Parse will parse a Paseto token and return it if it is valid.
func (x *PasetoMaker) Parse(t string) (*paseto.Token, error) {
	parsed, err := x.parser.ParseV4Local(x.symmetricKey, t, nil)
	if err != nil {
		return nil, fmt.Errorf("token: parsing token, %w", err)
	}
	return parsed, nil
}

func (x *PasetoMaker) RefreshTokenExpiration() time.Time {
	return time.Now().Add(x.refreshDur)
}

func (x *PasetoMaker) AccessTokenExpiration() time.Time {
	return time.Now().Add(x.accessDur)
}
