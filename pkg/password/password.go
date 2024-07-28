package password

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	// MaxBytes defines the maximum amount of bytes for a valid password.
	// It should always be 72 to satisfy BCRYPT.
	// Make sure to compare this with len().
	MaxBytes = 72

	// MinChars defines the minimum amount of utf8 runes for a valid password.
	// Make sure to compare this with [utf8.RuneCountInString()].
	MinChars = 8
)

// SafeString defines a safe password string.
// It cannot exceed [MaxBytes] or [MaxChars] and can not fall short of [MinChars].
// Its [Value()] will hash the underlying string with bcrypt.
// Its [String()] and [LogValue()] will mask the underlying string.
type SafeString string

type TooShortError struct{}

func (x TooShortError) Error() string {
	return fmt.Sprintf("password: must be at least %d characters long", MinChars)
}

type TooLongError struct {
	displayedForUser bool
}

func (x TooLongError) Error() string {
	// Don't reveal the hashing algo to user.
	if x.displayedForUser {
		return "password is too long"
	}
	return fmt.Sprintf("password: must be at most %d bytes long", MaxBytes)
}

// FromString will attempt to create a [SafeString] from a string.
// It will make sure the password is at least [MinChars] and at most [MaxBytes] long.
// Returns either [TooLongError] or [TooShortError] on error.
func FromString(s string) (SafeString, error) {
	if utf8.RuneCountInString(s) < MinChars {
		return "", TooShortError{}
	}
	if len(s) > MaxBytes {
		return "", TooLongError{displayedForUser: true}
	}

	// Password must contain at least one uppercase letter,
	// one lowercase letter and one digit.
	if !regexp.MustCompile(`[a-z]`).MatchString(s) ||
		!regexp.MustCompile(`[A-Z]`).MatchString(s) ||
		!regexp.MustCompile(`[0-9]`).MatchString(s) {
		return "", errors.New("password: must contain an uppercase and lowercase letter and a digit")
	}

	return SafeString(s), nil
}

// String will mask the underlying password string.
func (x SafeString) String() string {
	return "REDACTED"
}

// LogValue will mask the underlying password string.
func (x SafeString) LogValue() slog.Value {
	return slog.StringValue("REDACTED")
}

// Scan implements the database/sql Scanner interface.
func (x *SafeString) Scan(src any) error {
	if src == nil {
		x = new(SafeString)
		*x = ""
		return nil
	}

	switch s := src.(type) {
	case string:
		*x = SafeString(s)
		return nil
	case []byte:
		*x = SafeString(s)
		return nil
	}

	return fmt.Errorf("unsupported type: %T", src)
}

// Value implements the Valuer interface.
// It will hash the password with BCRYPT.
func (x SafeString) Value() (driver.Value, error) {
	if x == "" {
		return nil, errors.New("password is empty")
	}
	if utf8.RuneCountInString(string(x)) < MinChars {
		return nil, TooShortError{}
	}
	if len(x) > MaxBytes {
		return nil, TooLongError{}
	}

	return bcrypt.GenerateFromPassword([]byte(x), bcrypt.DefaultCost)
}
