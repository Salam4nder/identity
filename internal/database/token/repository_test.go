//go:build testdb
// +build testdb

package token_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/Salam4nder/identity/internal/database"
	"github.com/Salam4nder/identity/internal/database/token"
	"github.com/Salam4nder/identity/pkg/random"
	"github.com/google/uuid"
)

func TestInsert(t *testing.T) {
	ctx := context.Background()
	db, cleanup := database.SetupTestConn(token.Tablename)
	t.Cleanup(cleanup)

	commonID := uuid.NewString()
	commonEmail := random.Email()

	t.Run("OK", func(t *testing.T) {
		b := sha256.Sum256([]byte(commonEmail))
		h := hex.EncodeToString(b[:])
		tt := fmt.Sprintf("%s/%s", h, commonID)

		if err := token.Insert(ctx, db, tt); err != nil {
			t.Error("expected no error")
		}
	})

	t.Run("duplicate entry returns error", func(t *testing.T) {
		b := sha256.Sum256([]byte(commonEmail))
		h := hex.EncodeToString(b[:])
		tt := fmt.Sprintf("%s/%s", h, commonID)

		err := token.Insert(ctx, db, tt)
		if err == nil {
			t.Error("expected error")
		}
		if !errors.As(err, &database.DuplicateEntryError{}) {
			t.Error("expected duplicate entry error")
		}
	})

	t.Run("empty string returns error", func(t *testing.T) {
		err := token.Insert(ctx, db, "")
		if err == nil {
			t.Error("expected error")
		}
		if !errors.As(err, &database.InputError{}) {
			t.Error("expected input entry error")
		}
	})
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	db, cleanup := database.SetupTestConn(token.Tablename)
	t.Cleanup(cleanup)

	t.Run("OK", func(t *testing.T) {
		b := sha256.Sum256([]byte(random.Email()))
		h := hex.EncodeToString(b[:])
		tt := fmt.Sprintf("%s/%s", h, uuid.NewString())

		if err := token.Insert(ctx, db, tt); err != nil {
			t.Error("expected no error")
		}

		got, err := token.Get(ctx, db, tt)
		if err != nil {
			t.Error("expected no error")
		}
		if got != tt {
			t.Error("expected got to be equal to token")
		}
	})

	t.Run("empty string returns error", func(t *testing.T) {
		got, err := token.Get(ctx, db, "")
		if got != "" {
			t.Error("expected empty string")
		}
		if err == nil {
			t.Error("expected error")
		}
		if !errors.As(err, &database.InputError{}) {
			t.Error("expected input entry error")
		}
	})

	t.Run("not found returns error", func(t *testing.T) {
		got, err := token.Get(ctx, db, random.String(100))
		if got != "" {
			t.Error("expected empty string")
		}
		if err == nil {
			t.Error("expected error")
		}
		if !errors.As(err, &database.NotFoundError{}) {
			t.Error("expected input entry error")
		}
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	db, cleanup := database.SetupTestConn(token.Tablename)
	t.Cleanup(cleanup)

	t.Run("OK", func(t *testing.T) {
		b := sha256.Sum256([]byte(random.Email()))
		h := hex.EncodeToString(b[:])
		tt := fmt.Sprintf("%s/%s", h, uuid.NewString())

		if err := token.Insert(ctx, db, tt); err != nil {
			t.Error("expected no error")
		}

		if err := token.Delete(ctx, db, tt); err != nil {
			t.Error("expected no error")
		}
	})

	t.Run("not found returns error", func(t *testing.T) {
		err := token.Delete(ctx, db, random.String(100))
		if err == nil {
			t.Error("expected error")
		}
		if !errors.As(err, &database.RowsAffectedError{}) {
			t.Error("expected input entry error")
		}
	})
}
