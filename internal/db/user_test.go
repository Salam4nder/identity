//go:build testdb
// +build testdb

package db

import (
	"strings"
	"testing"
	"time"

	"github.com/Salam4nder/user/pkg/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSQL_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		params      CreateUserParams
		wantErr     bool
		requiredErr error
		preInsert   func() error
	}{
		{
			name: "Success",
			params: CreateUserParams{
				FullName:  util.RandomString(10),
				Email:     util.RandomEmail(),
				Password:  util.RandomString(10),
				CreatedAt: time.Now().UTC(),
			},
		},
		{
			name: "Duplicate email returns error",
			params: CreateUserParams{
				FullName:  "Kam Gam",
				Email:     "email@test.com",
				Password:  "password",
				CreatedAt: time.Now().UTC(),
			},
			preInsert: func() error {
				_, err := TestSQLConnPool.db.ExecContext(
					ctx,
					`INSERT INTO users (
                       full_name,
                       email,
                       password_hash,
                       created_at
                   ) VALUES (
                       'Kam Gam',
                       'email@test.com',
                       'password',
                       timezone('utc', now())
                       )`,
				)
				return err
			},
			wantErr:     true,
			requiredErr: ErrDuplicateEmail,
		},
		{
			name: "Name exceeds 255 chars returns err",
			params: CreateUserParams{
				FullName:  strings.Repeat("a", 256),
				Email:     util.RandomEmail(),
				Password:  util.RandomString(10),
				CreatedAt: time.Now().UTC(),
			},
			wantErr:     true,
			requiredErr: ErrStringTooLong,
		},
		{
			name: "Email exceeds 255 chars returns err",
			params: CreateUserParams{
				FullName:  util.RandomString(10),
				Email:     strings.Repeat("a", 256),
				Password:  util.RandomString(10),
				CreatedAt: time.Now().UTC(),
			},
			wantErr:     true,
			requiredErr: ErrStringTooLong,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preInsert != nil {
				require.NoError(t, test.preInsert())
			}

			got, err := TestSQLConnPool.CreateUser(ctx, test.params)
			if test.wantErr {
				require.Error(t, err)
				if test.requiredErr != nil {
					require.ErrorIs(t, err, test.requiredErr)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotEmpty(t, got.ID)
				require.NotEqual(t, got.ID, uuid.Nil)
				require.Equal(t, test.params.FullName, got.FullName)
				require.Equal(t, test.params.Email, got.Email)
				require.Equal(t, test.params.CreatedAt, got.CreatedAt)
				require.NotEmpty(t, got.PasswordHash)

				t.Cleanup(func() {
					_, err := TestSQLConnPool.db.ExecContext(
						ctx,
						`DELETE FROM users WHERE id = $1`,
						got.ID,
					)
					require.NoError(t, err)
				})
			}
		})
	}
}

func TestSQL_ReadUser(t *testing.T) {
	user, err := TestSQLConnPool.CreateUser(ctx, CreateUserParams{
		FullName:  util.RandomString(10),
		Email:     util.RandomEmail(),
		Password:  util.RandomString(10),
		CreatedAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	got, err := TestSQLConnPool.ReadUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, user.ID, got.ID)
	require.Equal(t, user.FullName, got.FullName)
	require.Equal(t, user.Email, got.Email)
	require.Equal(t, user.CreatedAt, got.CreatedAt)
	require.Equal(t, user.PasswordHash, got.PasswordHash)
	require.Equal(t, user.CreatedAt, got.CreatedAt)

	t.Cleanup(func() {
		_, err := TestSQLConnPool.db.ExecContext(
			ctx,
			`DELETE FROM users WHERE id = $1`,
			user.ID,
		)
		require.NoError(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		_, err := TestSQLConnPool.ReadUser(ctx, uuid.New())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("UUID is nil", func(t *testing.T) {
		_, err := TestSQLConnPool.ReadUser(ctx, uuid.Nil)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})
}
