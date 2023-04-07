package token

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payload contains the payload data of the token.
type Payload struct {
	ID        primitive.ObjectID `json:"id"`
	Email     string             `json:"username"`
	IssuedAt  time.Time          `json:"issued_at"`
	ExpiredAt time.Time          `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific email and duration.
func NewPayload(email string, duration time.Duration) (*Payload, error) {
	tokenID := primitive.NewObjectID()

	payload := &Payload{
		ID:        tokenID,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
