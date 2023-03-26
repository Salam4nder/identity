package storage

import (
	"testing"

	"github.com/Salam4nder/user/pkg/util"

	"github.com/stretchr/testify/assert"
)

func ValidateInsertOneParam(t *testing.T) {
	tests := []struct {
		name    string
		param   InsertParam
		wantErr bool
	}{
		{
			name: "valid param returns no error",
			param: InsertParam{
				FullName: util.RandomFullName(),
				Email:    util.RandomEmail(),
				Password: "123456",
			},
		},
		{
			name: "empty full name returns error",
			param: InsertParam{
				FullName: "",
				Email:    util.RandomEmail(),
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "empty email returns error",
			param: InsertParam{
				FullName: util.RandomFullName(),
				Email:    "",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "empty password returns error",
			param: InsertParam{
				FullName: util.RandomFullName(),
				Email:    util.RandomEmail(),
				Password: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.param.Validate()

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateParam(t *testing.T) {
	tests := []struct {
		name    string
		param   UpdateParam
		wantErr bool
	}{
		{
			name: "valid param 1 returns no error",
			param: UpdateParam{
				ID:       "456",
				FullName: util.RandomFullName(),
			},
		},
		{
			name: "valid param 2 returns no error",
			param: UpdateParam{
				ID:    "123",
				Email: util.RandomEmail(),
			},
		},
		{
			name: "empty ID returns error",
			param: UpdateParam{
				ID:    "",
				Email: util.RandomEmail(),
			},
			wantErr: true,
		},
		{
			name: "empty update fields returns error",
			param: UpdateParam{
				ID: "123",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.param.Validate()

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestValidateFilter(t *testing.T) {
	tests := []struct {
		name    string
		filter  Filter
		wantErr bool
	}{
		{
			name: "valid filter 1 returns no error",
			filter: Filter{
				FullName: util.RandomFullName(),
			},
		},
		{
			name: "valid filter 2 returns no error",
			filter: Filter{
				Email: util.RandomEmail(),
			},
		},
		{
			name: "valid filter 3 returns no error",
			filter: Filter{
				CreatedAt: util.RandomDate(),
			},
		},
		{
			name:    "invalid filter returns error",
			filter:  Filter{},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.filter.Validate()

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
