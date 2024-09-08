//go:build testdb
// +build testdb

package personalnumber_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Salam4nder/identity/internal/database"
	"github.com/Salam4nder/identity/internal/database/personalnumber"
)

var (
	db      *sql.DB
	cleanup func()
)

func init() {
	db, cleanup = database.SetupTestConn(personalnumber.Tablename)
}

func TestInsert(t *testing.T) {
	t.Cleanup(cleanup)

	n := uint64(4865998752658465)
	t.Run("OK", func(t *testing.T) {
		if err := personalnumber.Insert(context.Background(), db, n); err != nil {
			t.Errorf("expected no err, got %s", err.Error())
		}
	})

	t.Run("duplicate entry", func(t *testing.T) {
		err := personalnumber.Insert(context.Background(), db, n)
		if err == nil {
			t.Error("expected err")
		}
		if !errors.As(err, &database.DuplicateEntryError{}) {
			t.Error("expected duplicate entry error")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Cleanup(cleanup)

	n := uint64(4865998752658465)
	t.Run("OK", func(t *testing.T) {
		if err := personalnumber.Insert(context.Background(), db, n); err != nil {
			t.Errorf("expected no err, got %s", err.Error())
		}

		if err := personalnumber.Delete(context.Background(), db, n); err != nil {
			t.Errorf("expected no err, got %s", err.Error())
		}
	})

	t.Run("not found returns RowsAffectedError", func(t *testing.T) {
		err := personalnumber.Delete(context.Background(), db, 5)
		if err == nil {
			t.Error("expected err")
		}
		if !errors.As(err, &database.RowsAffectedError{}) {
			t.Error("expected rows affected error")
		}
	})
}
