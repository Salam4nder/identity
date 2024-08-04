//go:build testdb
// +build testdb

package db

import (
	"strings"
	"testing"
	"time"

	"github.com/Salam4nder/user/pkg/password"
	"github.com/Salam4nder/user/pkg/random"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestInsert(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("credentials")
	t.Cleanup(cleanup)

	randomParams := InsertParams{
		ID:        uuid.New(),
		Email:     random.Email(),
		Password:  password.SafeString(random.String(10)),
		CreatedAt: time.Now().UTC(),
	}

	t.Run("ok", func(t *testing.T) {
		t.Cleanup(cleanup)

		err := Insert(ctx, db, randomParams)
		require.NoError(t, err)

		got, err := Read(ctx, db, randomParams.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, randomParams.ID, got.ID)
		require.Equal(t, randomParams.Email, got.Email)
		require.NotEqual(t, randomParams.Password, got.PasswordHash)
		require.True(t, time.Now().After(got.CreatedAt))

		require.NoError(
			t,
			bcrypt.CompareHashAndPassword([]byte(got.PasswordHash), []byte(randomParams.Password)),
		)
	})

	t.Run("email exceeds 255 chars returns err", func(t *testing.T) {
		t.Cleanup(cleanup)

		randomParams.Email = strings.Repeat("a", 256)

		err := Insert(ctx, db, randomParams)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrStringTooLong)
		randomParams.Email = random.Email()
	})

	t.Run("duplicate email returns error", func(t *testing.T) {
		t.Cleanup(cleanup)

		err := Insert(ctx, db, InsertParams{
			ID:        uuid.New(),
			Email:     "email@email.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		})
		require.NoError(t, err)

		err = Insert(ctx, db, InsertParams{
			ID:        uuid.New(),
			Email:     "email@email.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrDuplicateEmail)
	})
}

func TestRead(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("credentials")
	t.Cleanup(cleanup)

	randomParams := InsertParams{
		ID:        uuid.New(),
		Email:     random.Email(),
		Password:  password.SafeString(random.String(10)),
		CreatedAt: time.Now().UTC(),
	}

	err := Insert(ctx, db, randomParams)
	require.NoError(t, err)

	got, err := Read(ctx, db, randomParams.ID)
	require.NoError(t, err)
	require.NotNil(t, got)

	t.Run("Not found", func(t *testing.T) {
		_, err := Read(ctx, db, uuid.New())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("InputError on nil UUID", func(t *testing.T) {
		_, err := Read(ctx, db, uuid.Nil)
		require.Error(t, err)
		require.ErrorAs(t, err, &InputError{})
	})
}

func TestReadByEmail(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("credentials")
	t.Cleanup(cleanup)

	randomParams := InsertParams{
		ID:        uuid.New(),
		Email:     random.Email(),
		Password:  password.SafeString(random.String(10)),
		CreatedAt: time.Now().UTC(),
	}

	err := Insert(ctx, db, randomParams)
	require.NoError(t, err)

	got, err := ReadByEmail(ctx, db, randomParams.Email)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, randomParams.ID, got.ID)
	require.Equal(t, randomParams.Email, got.Email)
	require.True(t, time.Now().After(got.CreatedAt))

	t.Run("Not found", func(t *testing.T) {
		_, err := ReadByEmail(ctx, db, random.Email())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("Email is empty", func(t *testing.T) {
		_, err := ReadByEmail(ctx, db, "")
		require.Error(t, err)
		require.ErrorAs(t, err, &InputError{})
	})

	t.Run("Email is not found", func(t *testing.T) {
		_, err := ReadByEmail(ctx, db, "ass@ass.com")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestUpdate(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("credentials")
	t.Cleanup(cleanup)

	randomParams := InsertParams{
		ID:        uuid.New(),
		Email:     random.Email(),
		Password:  password.SafeString(random.String(10)),
		CreatedAt: time.Now().UTC(),
	}

	t.Run("ok", func(t *testing.T) {
		t.Cleanup(cleanup)

		newEmail := "new@email.com"

		err := Insert(ctx, db, randomParams)
		require.NoError(t, err)

		err = Update(ctx, db, UpdateParams{
			ID:    randomParams.ID,
			Email: newEmail,
		})
		require.NoError(t, err)

		got, err := Read(ctx, db, randomParams.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, randomParams.ID, got.ID)
		require.Equal(t, newEmail, got.Email)
		require.True(t, time.Now().After(got.CreatedAt))
	})

	t.Run("email exceeds 255 chars returns err", func(t *testing.T) {
		t.Cleanup(cleanup)

		ID := uuid.New()

		err := Insert(ctx, db, InsertParams{
			ID:        ID,
			Email:     random.Email(),
			Password:  password.SafeString(random.String(10)),
			CreatedAt: time.Now().UTC(),
		})
		require.NoError(t, err)

		err = Update(ctx, db, UpdateParams{
			ID:    ID,
			Email: strings.Repeat("a", 256),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrStringTooLong)
	})
}

func TestDelete(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("credentials")
	t.Cleanup(cleanup)

	ID := uuid.New()

	err := Insert(ctx, db, InsertParams{
		ID:        ID,
		Email:     random.Email(),
		Password:  password.SafeString(random.String(15)),
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	err = Delete(ctx, db, ID)
	require.NoError(t, err)

	t.Run("Not found", func(t *testing.T) {
		err := Delete(ctx, db, ID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNoRowsAffected)
	})
}
