package db

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// Common db errors.
var (
	ErrDuplicateEmail = errors.New("db: duplicate email")
	ErrStringTooLong  = errors.New("db: string too long")
	ErrUserNotFound   = errors.New("db: user not found")
)

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
