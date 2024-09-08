//go:build testdb
// +build testdb

package credentials_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/Salam4nder/identity/internal/database"
	"github.com/Salam4nder/identity/internal/database/credentials"
	"github.com/Salam4nder/identity/pkg/password"
	"github.com/Salam4nder/identity/pkg/random"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	db      *sql.DB
	cleanup func()
)

func init() {
	db, cleanup = database.SetupTestConn(credentials.Tablename)
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(cleanup)

	randomParams := credentials.InsertParams{
		ID:        uuid.New(),
		Email:     random.Email(),
		Password:  password.SafeString(random.String(10)),
		CreatedAt: time.Now().UTC(),
	}

	t.Run("ok", func(t *testing.T) {
		t.Cleanup(cleanup)

		err := credentials.Insert(ctx, db, randomParams)
		require.NoError(t, err)

		got, err := credentials.Read(ctx, db, randomParams.ID)
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

		err := credentials.Insert(ctx, db, randomParams)
		require.Error(t, err)
		randomParams.Email = random.Email()
	})

	t.Run("duplicate email returns error", func(t *testing.T) {
		t.Cleanup(cleanup)

		err := credentials.Insert(ctx, db, credentials.InsertParams{
			ID:        uuid.New(),
			Email:     "email@email.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		})
		require.NoError(t, err)

		err = credentials.Insert(ctx, db, credentials.InsertParams{
			ID:        uuid.New(),
			Email:     "email@email.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		})
		require.Error(t, err)
		require.ErrorAs(t, err, &database.DuplicateEntryError{})
	})
}

func TestRead(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(cleanup)

	randomParams := credentials.InsertParams{
		ID:        uuid.New(),
		Email:     random.Email(),
		Password:  password.SafeString(random.String(10)),
		CreatedAt: time.Now().UTC(),
	}

	err := credentials.Insert(ctx, db, randomParams)
	require.NoError(t, err)

	got, err := credentials.Read(ctx, db, randomParams.ID)
	require.NoError(t, err)
	require.NotNil(t, got)

	t.Run("Not found", func(t *testing.T) {
		_, err := credentials.Read(ctx, db, uuid.New())
		require.Error(t, err)
		require.ErrorAs(t, err, &database.NotFoundError{})
	})

	t.Run("InputError on nil UUID", func(t *testing.T) {
		_, err := credentials.Read(ctx, db, uuid.Nil)
		require.Error(t, err)
		require.ErrorAs(t, err, &database.InputError{})
	})
}

func TestReadByEmail(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(cleanup)

	randomParams := credentials.InsertParams{
		ID:        uuid.New(),
		Email:     random.Email(),
		Password:  password.SafeString(random.String(10)),
		CreatedAt: time.Now().UTC(),
	}

	err := credentials.Insert(ctx, db, randomParams)
	require.NoError(t, err)

	got, err := credentials.ReadByEmail(ctx, db, randomParams.Email)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, randomParams.ID, got.ID)
	require.Equal(t, randomParams.Email, got.Email)
	require.True(t, time.Now().After(got.CreatedAt))

	t.Run("Not found", func(t *testing.T) {
		_, err := credentials.ReadByEmail(ctx, db, random.Email())
		require.Error(t, err)
		require.ErrorAs(t, err, &database.NotFoundError{})
	})

	t.Run("Email is empty", func(t *testing.T) {
		_, err := credentials.ReadByEmail(ctx, db, "")
		require.Error(t, err)
		require.ErrorAs(t, err, &database.InputError{})
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(cleanup)

	randomParams := credentials.InsertParams{
		ID:        uuid.New(),
		Email:     random.Email(),
		Password:  password.SafeString(random.String(10)),
		CreatedAt: time.Now().UTC(),
	}

	t.Run("OK", func(t *testing.T) {
		t.Cleanup(cleanup)

		newEmail := "new@email.com"

		err := credentials.Insert(ctx, db, randomParams)
		require.NoError(t, err)

		err = credentials.Update(ctx, db, credentials.UpdateParams{
			ID:    randomParams.ID,
			Email: newEmail,
		})
		require.NoError(t, err)

		got, err := credentials.Read(ctx, db, randomParams.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, randomParams.ID, got.ID)
		require.Equal(t, newEmail, got.Email)
		require.True(t, time.Now().After(got.CreatedAt))
	})

	t.Run("email exceeds 255 chars returns err", func(t *testing.T) {
		t.Cleanup(cleanup)

		ID := uuid.New()

		err := credentials.Insert(ctx, db, credentials.InsertParams{
			ID:        ID,
			Email:     random.Email(),
			Password:  password.SafeString(random.String(10)),
			CreatedAt: time.Now().UTC(),
		})
		require.NoError(t, err)

		err = credentials.Update(ctx, db, credentials.UpdateParams{
			ID:    ID,
			Email: strings.Repeat("a", 256),
		})
		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		err := credentials.Update(ctx, db, credentials.UpdateParams{
			ID:    uuid.New(),
			Email: strings.Repeat("a", 23),
		})
		require.Error(t, err)
		require.ErrorAs(t, err, &database.RowsAffectedError{})
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(cleanup)

	ID := uuid.New()

	t.Run("OK", func(t *testing.T) {
		err := credentials.Insert(ctx, db, credentials.InsertParams{
			ID:        ID,
			Email:     random.Email(),
			Password:  password.SafeString(random.String(15)),
			CreatedAt: time.Now(),
		})
		require.NoError(t, err)

		err = credentials.Delete(ctx, db, ID)
		require.NoError(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		err := credentials.Delete(ctx, db, ID)
		require.Error(t, err)
		require.ErrorAs(t, err, &database.RowsAffectedError{})
	})
}

func TestVerify(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(cleanup)

	ID := uuid.New()

	t.Run("OK", func(t *testing.T) {
		err := credentials.Insert(ctx, db, credentials.InsertParams{
			ID:        ID,
			Email:     random.Email(),
			Password:  password.SafeString(random.String(15)),
			CreatedAt: time.Now(),
		})
		require.NoError(t, err)

		err = credentials.Verify(ctx, db, ID)
		require.NoError(t, err)

		cred, err := credentials.Read(ctx, db, ID)
		require.NoError(t, err)
		require.NotNil(t, cred.VerifiedAt)
	})
}
