package db

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// Common db errors.
var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrUserNotFound   = errors.New("user not found")
)

// IsSentinelErr is a sentinel error.
func IsSentinelErr(err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}

	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		return true
	}

	return false
}

// SentinelErr returns a PSQL sentinel error.
func SentinelErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	var pqErr *pq.Error

	// nolint: errorlint
	pqErr = err.(*pq.Error)

	switch pqErr.Code.Name() {
	case "unique_violation":
		return ErrDuplicateEmail

	// TODO: Add more cases here.
	default:
		return err
	}
}
