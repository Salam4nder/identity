//go:build testdb
// +build testdb

package db

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Salam4nder/user/pkg/random"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSQL_CreateUser(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("users")
	t.Cleanup(cleanup)

	t.Run("ok", func(t *testing.T) {
		t.Cleanup(cleanup)

		ID := uuid.New()
		fullName := random.FullName()
		email := random.Email()
		password := random.String(10)
		createdAt := time.Now().UTC()

		err := db.CreateUser(ctx, CreateUserParams{
			ID:        ID,
			FullName:  fullName,
			Email:     email,
			Password:  password,
			CreatedAt: createdAt,
		})
		require.NoError(t, err)

		got, err := db.ReadUser(ctx, ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, ID, got.ID)
		require.Equal(t, fullName, got.FullName)
		require.Equal(t, email, got.Email)
		require.Equal(t, createdAt, got.CreatedAt)
	})

	t.Run("name exceeds 255 chars returns err", func(t *testing.T) {
		t.Cleanup(cleanup)

		err := db.CreateUser(ctx, CreateUserParams{
			ID:        uuid.New(),
			FullName:  strings.Repeat("a", 256),
			Email:     random.Email(),
			Password:  random.String(10),
			CreatedAt: time.Now().UTC(),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrStringTooLong)
	})

	t.Run("email exceeds 255 chars returns err", func(t *testing.T) {
		t.Cleanup(cleanup)

		err := db.CreateUser(ctx, CreateUserParams{
			ID:        uuid.New(),
			FullName:  random.FullName(),
			Email:     strings.Repeat("a", 256),
			Password:  random.String(10),
			CreatedAt: time.Now().UTC(),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrStringTooLong)
	})

	t.Run("duplicate email returns error", func(t *testing.T) {
		t.Cleanup(cleanup)

		err := db.CreateUser(ctx, CreateUserParams{
			ID:        uuid.New(),
			FullName:  "Kam Gam",
			Email:     "email@email.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		})
		require.NoError(t, err)

		err = db.CreateUser(ctx, CreateUserParams{
			ID:        uuid.New(),
			FullName:  "Kam Gam",
			Email:     "email@email.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrDuplicateEmail)
	})
}

func TestSQL_ReadUser(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("users")
	t.Cleanup(cleanup)

	ID := uuid.New()
	fullName := random.FullName()
	email := random.Email()
	password := random.String(10)
	createdAt := time.Now().UTC()

	err := db.CreateUser(ctx, CreateUserParams{
		ID:        ID,
		FullName:  fullName,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	})
	require.NoError(t, err)

	ctx := context.Background()
	got, err := db.ReadUser(ctx, ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, ID, got.ID)
	require.Equal(t, fullName, got.FullName)
	require.Equal(t, email, got.Email)
	require.Equal(t, createdAt, got.CreatedAt)

	t.Run("Not found", func(t *testing.T) {
		_, err := db.ReadUser(ctx, uuid.New())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("UUID is nil", func(t *testing.T) {
		_, err := db.ReadUser(ctx, uuid.Nil)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestSQL_ReadUserByEmail(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("users")
	t.Cleanup(cleanup)

	ID := uuid.New()
	fullName := random.FullName()
	email := random.Email()
	password := random.String(10)
	createdAt := time.Now().UTC()

	err := db.CreateUser(ctx, CreateUserParams{
		ID:        ID,
		FullName:  fullName,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	})
	require.NoError(t, err)

	ctx := context.Background()
	got, err := db.ReadUserByEmail(ctx, email)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, ID, got.ID)
	require.Equal(t, fullName, got.FullName)
	require.Equal(t, email, got.Email)
	require.Equal(t, createdAt, got.CreatedAt)
	require.Equal(t, createdAt, got.CreatedAt)

	t.Run("Not found", func(t *testing.T) {
		_, err := db.ReadUserByEmail(ctx, random.Email())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("Email is empty", func(t *testing.T) {
		_, err := db.ReadUserByEmail(ctx, "")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestSQL_UpdateUser(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("users")
	t.Cleanup(cleanup)

	t.Run("ok", func(t *testing.T) {
		t.Cleanup(cleanup)

		ID := uuid.New()
		fullName := random.FullName()
		email := random.Email()
		password := random.String(10)
		createdAt := time.Now().UTC()

		newFullName := "New Name"
		newEmail := "new@email.com"

		err := db.CreateUser(ctx, CreateUserParams{
			ID:        ID,
			FullName:  fullName,
			Email:     email,
			Password:  password,
			CreatedAt: createdAt,
		})
		require.NoError(t, err)

		ctx := context.Background()
		err = db.UpdateUser(ctx, UpdateUserParams{
			ID:       ID,
			FullName: newFullName,
			Email:    newEmail,
		})
		require.NoError(t, err)

		got, err := db.ReadUser(ctx, ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, ID, got.ID)
		require.Equal(t, newFullName, got.FullName)
		require.Equal(t, newEmail, got.Email)
		require.Equal(t, createdAt, got.CreatedAt)
	})

	t.Run("name exceeds 255 chars returns err", func(t *testing.T) {
		t.Cleanup(cleanup)

		ID := uuid.New()

		err := db.CreateUser(ctx, CreateUserParams{
			ID:        uuid.New(),
			FullName:  random.FullName(),
			Email:     random.Email(),
			Password:  random.String(10),
			CreatedAt: time.Now().UTC(),
		})
		require.NoError(t, err)

		ctx := context.Background()
		err = db.UpdateUser(ctx, UpdateUserParams{
			ID:       ID,
			FullName: strings.Repeat("a", 256),
			Email:    random.Email(),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrStringTooLong)
	})

	t.Run("email exceeds 255 chars returns err", func(t *testing.T) {
		t.Cleanup(cleanup)

		ID := uuid.New()

		err := db.CreateUser(ctx, CreateUserParams{
			ID:        uuid.New(),
			FullName:  random.FullName(),
			Email:     random.Email(),
			Password:  random.String(10),
			CreatedAt: time.Now().UTC(),
		})
		require.NoError(t, err)

		ctx := context.Background()
		err = db.UpdateUser(ctx, UpdateUserParams{
			ID:       ID,
			FullName: random.FullName(),
			Email:    strings.Repeat("a", 256),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrStringTooLong)
	})
}

func TestSQL_DeleteUser(t *testing.T) {
	db, cleanup := NewTestSQLConnPool("users")
	t.Cleanup(cleanup)

	ID := uuid.New()

	err := db.CreateUser(ctx, CreateUserParams{
		ID:        ID,
		FullName:  random.FullName(),
		Email:     random.Email(),
		Password:  random.String(10),
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	ctx := context.Background()
	err = db.DeleteUser(ctx, ID)
	require.NoError(t, err)

	t.Run("Not found", func(t *testing.T) {
		err := db.DeleteUser(ctx, ID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNoRowsAffected)
	})
}
