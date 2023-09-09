package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// pasetoMaker is a PASETO token maker.
type pasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker.
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf(
			"invalid key size: must be exactly %d characters",
			chacha20poly1305.KeySize)
	}

	maker := &pasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// NewToken creates a new token for a specific email and duration.
func (x *pasetoMaker) NewToken(
	email string,
	duration time.Duration,
) (string, *Payload, error) {
	payload := NewPayload(email, duration)

	token, err := x.paseto.Encrypt(x.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken checks if the token is valid.
func (x *pasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	if err := x.paseto.Decrypt(
		token,
		x.symmetricKey,
		payload,
		nil,
	); err != nil {
		return nil, ErrInvalidToken
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
