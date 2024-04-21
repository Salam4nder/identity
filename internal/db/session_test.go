package db

import (
	"context"
	"testing"
	"time"

	"github.com/Salam4nder/user/pkg/random"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	db, cleanup := NewTestSQLConnPool()
	t.Cleanup(cleanup)

	t.Run("ok", func(t *testing.T) {
		err := db.CreateSession(context.Background(), CreateSessionParams{
			ID:           uuid.New(),
			Email:        random.Email(),
			ClientIP:     random.String(15),
			UserAgent:    random.String(20),
			RefreshToken: random.String(32),
			ExpiresAt:    time.Now().Add(time.Hour),
		})
		require.NoError(t, err)
	})

	t.Run("missing fields", func(t *testing.T) {
		err := db.CreateSession(context.Background(), CreateSessionParams{})
		require.Error(t, err)
	})
}
