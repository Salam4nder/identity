package token

import (
	"time"
)

// Maker is an interface for managing tokens
type Maker interface {
	NewToken(email string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
