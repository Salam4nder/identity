package user

import (
	"context"
	"testing"

	"github.com/Salam4nder/user/pkg/util"
	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestInsertOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success returns id as string and nil", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		params := CreateParam{
			FullName:  util.RandomFullName(),
			Email:     util.RandomEmail(),
			Password:  util.RandomString(10),
			CreatedAt: util.RandomDate(),
		}

		service := NewService(mt.Coll)

		id, err := service.InsertOne(context.TODO(), params)

		assert.Nil(t, err)
		assert.NotEmpty(t, id)
		assert.IsType(t, id, "")
	})

	mt.Run("error returns empty string and err", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    11000,
				Message: "duplicate key error",
			}))

		params := CreateParam{
			FullName:  util.RandomFullName(),
			Email:     util.RandomEmail(),
			Password:  util.RandomString(10),
			CreatedAt: util.RandomDate(),
		}

		service := NewService(mt.Coll)

		id, err := service.InsertOne(context.TODO(), params)

		assert.NotNil(t, err)
		assert.Empty(t, id)
		assert.IsType(t, id, "")
	})
}
