package user

import (
	"testing"

	"github.com/Salam4nder/user/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestValidateFilter(t *testing.T) {
	tests := []struct {
		name    string
		filter  Filter
		wantErr bool
	}{
		{
			name: "valid filter 1 returns no error",
			filter: Filter{
				FullName: "John",
			},
		},
		{
			name: "valid filter 2 returns no error",
			filter: Filter{
				Email: "lma@mail.com",
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
