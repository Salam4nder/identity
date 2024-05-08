package token

import (
	"time"

	"github.com/google/uuid"
)

// Payload contains the payload data of the token.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific email and duration.
func NewPayload(email string, duration time.Duration) *Payload {
	return &Payload{
		ID:        uuid.New(),
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
}

// Valid checks if the token payload is valid.
func (x *Payload) Valid() error {
	if time.Now().After(x.ExpiresAt) {
		return ErrExpiredToken
	}

	return nil
}
