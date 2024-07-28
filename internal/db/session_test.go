//go:build testdb
// +build testdb

package db

import (
	"context"
	"testing"
	"time"

	"github.com/Salam4nder/user/pkg/password"
	"github.com/Salam4nder/user/pkg/random"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	dbSess, sessCleanup := NewTestSQLConnPool("sessions")
	t.Cleanup(sessCleanup)

	dbUser, userCleanup := NewTestSQLConnPool("users")
	t.Cleanup(userCleanup)

	t.Run("ok", func(t *testing.T) {
		email := random.Email()
		err := dbUser.CreateUser(context.Background(), CreateUserParams{
			ID:       uuid.New(),
			Email:    email,
			Password: password.SafeString(random.String(32)),
		})
		require.NoError(t, err)

		err = dbSess.CreateSession(context.Background(), CreateSessionParams{
			ID:           uuid.New(),
			Email:        email,
			ClientIP:     random.String(15),
			UserAgent:    random.String(20),
			RefreshToken: random.String(32),
			ExpiresAt:    time.Now().Add(time.Hour),
		})
		require.NoError(t, err)

		t.Cleanup(sessCleanup)
		t.Cleanup(userCleanup)
	})

	t.Run("missing fields", func(t *testing.T) {
		email := random.Email()
		err := dbUser.CreateUser(context.Background(), CreateUserParams{
			ID:       uuid.New(),
			Email:    email,
			Password: password.SafeString(random.String(32)),
		})
		require.NoError(t, err)

		err = dbSess.CreateSession(context.Background(), CreateSessionParams{})
		require.Error(t, err)
	})
}
