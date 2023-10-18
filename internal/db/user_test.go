//go:build testdb
// +build testdb

package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQL_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		params  CreateUserParams
		wantErr bool
	}{
		{
			name: "Success",
			params: CreateUserParams{
				FullName:  "Kam Gam",
				Email:     "email@test.com",
				Password:  "password",
				CreatedAt: time.Now(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := TestSQLConnPool.CreateUser(ctx, test.params)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotEmpty(t, got.ID)
				require.Equal(t, test.params.FullName, got.FullName)
				require.NotEmpty(t, got.PasswordHash)
			}
		})
	}
}
