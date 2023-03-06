package user

import (
	"context"
	"testing"

	"github.com/Salam4nder/user/pkg/util"
	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestFindOneByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success returns user and nil", func(mt *mtest.T) {
		idObj := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1, "foo.bar", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: idObj},
					{Key: "full_name", Value: "John Doe"},
					{Key: "email", Value: "lmao@gmail.com"},
				}))

		service := NewService(mt.Coll)

		user, err := service.FindOneByID(context.TODO(), idObj.Hex())

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.IsType(t, user, User{})
		assert.Equal(t, user.ID, idObj)
		assert.Equal(t, user.FullName, "John Doe")
		assert.Equal(t, user.Email, "lmao@gmail.com")
	})

	mt.Run("error returns empty user and err", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse())

		service := NewService(mt.Coll)

		user, err := service.FindOneByID(context.TODO(), "123")

		assert.NotNil(t, err)
		assert.Empty(t, user)
		assert.IsType(t, user, User{})
	})
}

func TestFindOneByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success returns user and nil", func(mt *mtest.T) {
		idObj := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1, "foo.bar", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: idObj},
					{Key: "full_name", Value: "John Doe"},
					{Key: "email", Value: "haha@gmail.com"},
				}))

		service := NewService(mt.Coll)

		user, err := service.FindOneByEmail(context.TODO(), "haha@gmail.com")

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.IsType(t, user, User{})
		assert.Equal(t, user.ID, idObj)
		assert.Equal(t, user.FullName, "John Doe")
		assert.Equal(t, user.Email, "haha@gmail.com")
	})

	mt.Run("error returns empty user and err", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse())

		service := NewService(mt.Coll)

		user, err := service.FindOneByEmail(context.TODO(), "")

		assert.NotNil(t, err)
		assert.Empty(t, user)
		assert.IsType(t, user, User{})
	})
}
