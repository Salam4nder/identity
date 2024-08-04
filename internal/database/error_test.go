package database_test

import (
	"errors"
	"testing"

	"github.com/Salam4nder/user/internal/database"
	"github.com/lib/pq"
)

const PQErrUniqueViolationCode = "23505"

func Test_IsPSQLDuplicateEntryError(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		err := &pq.Error{
			Code: PQErrUniqueViolationCode,
		}

		if !database.IsPSQLDuplicateEntryError(err) {
			t.Error("expected true")
		}

	})

	t.Run("should be false", func(t *testing.T) {
		err := errors.New("ass")

		if database.IsPSQLDuplicateEntryError(err) {
			t.Error("expected false")
		}
	})
}
