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
				FullName:  random.FullName(),
				Email:     random.Email(),
				Password:  random.String(10),
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
				Email:     random.Email(),
				Password:  random.String(10),
				CreatedAt: time.Now().UTC(),
			},
			wantErr:     true,
			requiredErr: ErrStringTooLong,
		},
		{
			name: "Email exceeds 255 chars returns err",
			params: CreateUserParams{
				FullName:  random.FullName(),
				Email:     strings.Repeat("a", 256),
				Password:  random.String(10),
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
			t.Cleanup(func() {
				_, err := TestSQLConnPool.db.ExecContext(
					ctx,
					`DELETE FROM users`,
				)
				require.NoError(t, err)
			})

			ctx := context.Background()
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

			}
		})
	}
}

func TestSQL_ReadUser(t *testing.T) {
	user, err := TestSQLConnPool.CreateUser(ctx, CreateUserParams{
		FullName:  random.FullName(),
		Email:     random.Email(),
		Password:  random.String(10),
		CreatedAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err := TestSQLConnPool.db.ExecContext(
			ctx,
			`DELETE FROM users`,
		)
		require.NoError(t, err)
	})

	ctx := context.Background()
	got, err := TestSQLConnPool.ReadUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, user.ID, got.ID)
	require.Equal(t, user.FullName, got.FullName)
	require.Equal(t, user.Email, got.Email)
	require.Equal(t, user.CreatedAt, got.CreatedAt)
	require.Equal(t, user.PasswordHash, got.PasswordHash)
	require.Equal(t, user.CreatedAt, got.CreatedAt)

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

func TestSQL_ReadUserByEmail(t *testing.T) {
	user, err := TestSQLConnPool.CreateUser(ctx, CreateUserParams{
		FullName:  random.FullName(),
		Email:     random.Email(),
		Password:  random.String(10),
		CreatedAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err := TestSQLConnPool.db.ExecContext(
			ctx,
			`DELETE FROM users`,
		)
		require.NoError(t, err)
	})

	ctx := context.Background()
	got, err := TestSQLConnPool.ReadUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, user.ID, got.ID)
	require.Equal(t, user.FullName, got.FullName)
	require.Equal(t, user.Email, got.Email)
	require.Equal(t, user.CreatedAt, got.CreatedAt)
	require.Equal(t, user.PasswordHash, got.PasswordHash)
	require.Equal(t, user.CreatedAt, got.CreatedAt)

	t.Run("Not found", func(t *testing.T) {
		_, err := TestSQLConnPool.ReadUserByEmail(ctx, random.Email())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("Email is empty", func(t *testing.T) {
		_, err := TestSQLConnPool.ReadUserByEmail(ctx, "")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestSQL_UpdateUser(t *testing.T) {
	type args struct {
		ctx    context.Context
		params UpdateUserParams
	}
	tests := []struct {
		name          string
		args          args
		preUpdateArgs CreateUserParams
		wantErr       bool
		requiredErr   error
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				params: UpdateUserParams{
					FullName: random.FullName(),
					Email:    random.Email(),
				},
			},
			preUpdateArgs: CreateUserParams{
				FullName:  random.FullName(),
				Email:     random.Email(),
				Password:  random.String(10),
				CreatedAt: time.Now().UTC(),
			},
		},
		{
			name: "Name exceeds 255 chars returns err",
			args: args{
				ctx: context.Background(),
				params: UpdateUserParams{
					FullName: strings.Repeat("a", 256),
				},
			},
			preUpdateArgs: CreateUserParams{
				FullName:  random.FullName(),
				Email:     random.Email(),
				Password:  random.String(10),
				CreatedAt: time.Now().UTC(),
			},
			wantErr:     true,
			requiredErr: ErrStringTooLong,
		},
		{
			name: "Email exceeds 255 chars returns err",
			args: args{
				ctx: context.Background(),
				params: UpdateUserParams{
					Email: strings.Repeat("a", 256),
				},
			},
			preUpdateArgs: CreateUserParams{
				FullName:  random.FullName(),
				Email:     random.Email(),
				Password:  random.String(10),
				CreatedAt: time.Now().UTC(),
			},
			wantErr:     true,
			requiredErr: ErrStringTooLong,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			preUpdateUser, err := TestSQLConnPool.CreateUser(ctx, test.preUpdateArgs)
			require.NoError(t, err)

			t.Cleanup(func() {
				_, err := TestSQLConnPool.db.ExecContext(
					ctx,
					`DELETE FROM users`,
				)
				require.NoError(t, err)
			})

			test.args.params.ID = preUpdateUser.ID
			got, err := TestSQLConnPool.UpdateUser(test.args.ctx, test.args.params)
			if test.wantErr {
				require.Error(t, err)
				if test.requiredErr != nil {
					require.ErrorIs(t, err, test.requiredErr)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, test.args.params.ID, got.ID)
				require.Equal(t, test.args.params.FullName, got.FullName)
				require.Equal(t, test.args.params.Email, got.Email)
				require.Equal(t, preUpdateUser.CreatedAt, got.CreatedAt)
				require.NotEmpty(t, got.UpdatedAt)
				require.NotEqual(t, preUpdateUser.UpdatedAt, got.UpdatedAt)
				require.True(t, preUpdateUser.CreatedAt.Before(*got.UpdatedAt))
			}
		})
	}
}

func TestSQL_DeleteUser(t *testing.T) {
	createdUser, err := TestSQLConnPool.CreateUser(ctx, CreateUserParams{
		FullName:  random.FullName(),
		Email:     random.Email(),
		Password:  random.String(10),
		CreatedAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, createdUser.ID)

	ctx := context.Background()
	err = TestSQLConnPool.DeleteUser(ctx, createdUser.ID)
	require.NoError(t, err)

	t.Run("Not found", func(t *testing.T) {
		err := TestSQLConnPool.DeleteUser(ctx, createdUser.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})
}
