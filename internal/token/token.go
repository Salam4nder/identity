package token

import (
	"log/slog"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/Salam4nder/identity/proto/gen"
)

// SafeString has it's common stringers masked.
// Should be created internally.
type SafeString string

// String will mask the underlying password string.
func (x SafeString) String() string {
	return "REDACTED"
}

// LogValue will mask the underlying password string.
func (x SafeString) LogValue() slog.Value {
	return slog.StringValue("REDACTED")
}

func fromString(s string) SafeString {
	return SafeString(s)
}

// Maker is an abstract interface for making and verifying access and refresh tokens.
type Maker interface {
	MakeAccessToken(identifier any, strategy gen.Strategy) (SafeString, error)
	MakeRefreshToken(identifier any, strategy gen.Strategy) (SafeString, error)
	Parse(token string) (*paseto.Token, error)
	RefreshTokenExpiration() time.Time
	AccessTokenExpiration() time.Time
}
