package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/lib/pq"
)

const PQErrUniqueViolationCode = "23505"

func Test_IsSentinelErr(t *testing.T) {
	t.Run("sql.ErrNoRows returns true", func(t *testing.T) {
		if !IsSentinelErr(sql.ErrNoRows) {
			t.Error("sql.ErrNoRows should be a sentinel error")
		}
	})

	t.Run("pq.Error returns true", func(t *testing.T) {
		err := &pq.Error{
			Code: "unique_violation",
		}

		if !IsSentinelErr(err) {
			t.Error("pq.Error should be a sentinel error")
		}
	})

	t.Run("non-sentinel error returns false", func(t *testing.T) {
		err := errors.New("some error")

		if IsSentinelErr(err) {
			t.Error("errors.New should not be a sentinel error")
		}
	})
}

func Test_SentinelErr(t *testing.T) {
	t.Run("sql.ErrNoRows returns ErrUserNotFound", func(t *testing.T) {
		if !errors.Is(SentinelErr(sql.ErrNoRows), ErrUserNotFound) {
			t.Error("sql.ErrNoRows should be ErrUserNotFound")
		}
	})

	t.Run("pq.Error returns ErrDuplicateEmail", func(t *testing.T) {
		err := &pq.Error{
			Code: PQErrUniqueViolationCode,
		}

		if !errors.Is(SentinelErr(err), ErrDuplicateEmail) {
			t.Error("pq.Error should be ErrDuplicateEmail")
		}
	})

	t.Run("non-specified sentinel error returns itself", func(t *testing.T) {
		err := &pq.Error{
			Code: "some_other_code",
		}

		if !errors.Is(SentinelErr(err), err) {
			t.Error("non-specified sentinel error should return itself")
		}
	})
}
