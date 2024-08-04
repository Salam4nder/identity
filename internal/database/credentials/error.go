package credentials

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrDuplicateEmail       = errors.New("db: duplicate email")
	ErrStringTooLong        = errors.New("db: string too long")
	ErrUserNotFound         = errors.New("db: user not found")
	ErrNoRowsAffected       = errors.New("db: no rows affected")
	ErrMultipleRowsAffected = errors.New("db: multiple rows affected")
)

// InputError is returned in case of input errors that are meant to be
// sent back to the user.
type InputError struct {
	Field string
	Value any
}

func (x InputError) Error() string {
	return fmt.Sprintf("input error on %s with value %v", x.Field, x.Value)
}

// IsSentinelErr is a sentinel error.
func IsSentinelErr(err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}

	var pqErr *pq.Error

	return errors.As(err, &pqErr)
}

// SentinelErr returns a PSQL sentinel error.
func SentinelErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	// nolint: errorlint
	pqErr := err.(*pq.Error)

	switch pqErr.Code.Name() {
	case "unique_violation":
		return ErrDuplicateEmail
	case "string_data_right_truncation":
		return ErrStringTooLong
	default:
		return err
	}
}
