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

		assert.Error(t, err)
		assert.NotNil(t, err)
		assert.Empty(t, user)
		assert.IsType(t, user, User{})
	})

	mt.Run("invalid id returns empty user and err", func(mt *mtest.T) {
		service := NewService(mt.Coll)

		user, err := service.FindOneByID(context.TODO(), "123")

		assert.Error(t, err)
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

func TestFindByFilter(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success returns list of users and nil", func(mt *mtest.T) {
		idObj := primitive.NewObjectID()
		idObj2 := primitive.NewObjectID()
		randomDate := util.RandomDate()

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: idObj},
			{Key: "full_name", Value: "John Doe"},
			{Key: "email", Value: "email@email.com"},
			{Key: "created_at", Value: randomDate},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: idObj2},
			{Key: "full_name", Value: "Mary Hopkins"},
			{Key: "email", Value: "email2@email.com"},
			{Key: "created_at", Value: randomDate},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		service := NewService(mt.Coll)

		filter := Filter{
			CreatedAt: randomDate,
		}

		users, err := service.FindByFilter(context.TODO(), filter)

		assert.Nil(t, err)
		assert.NotNil(t, users)
		assert.IsType(t, users, []User{})
		assert.Equal(t, users[0].ID, idObj)
		assert.Equal(t, users[0].FullName, "John Doe")
		assert.Equal(t, users[0].Email, "email@email.com")
		assert.Equal(t, users[0].CreatedAt, randomDate)
		assert.Equal(t, users[1].ID, idObj2)
		assert.Equal(t, users[1].FullName, "Mary Hopkins")
		assert.Equal(t, users[1].Email, "email2@email.com")
	})

	mt.Run("encoding error returns empty list and err", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse())

		service := NewService(mt.Coll)

		filter := Filter{
			CreatedAt: util.RandomDate(),
		}

		users, err := service.FindByFilter(context.TODO(), filter)

		assert.Error(t, err)
		assert.NotNil(t, err)
		assert.Empty(t, users)
		assert.IsType(t, users, []User{})
	})
}

func TestUpdateOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success returns updated user and nil", func(mt *mtest.T) {
		idObj := primitive.NewObjectID()
		randomDate := util.RandomDate()

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: idObj},
				{Key: "full_name", Value: "John Doe"},
				{Key: "email", Value: "email@veryhotmale.com"},
				{Key: "created_at", Value: randomDate},
			}},
		})

		service := NewService(mt.Coll)

		updateParam := UpdateParam{
			ID:       idObj.Hex(),
			FullName: "John Doe",
			Email:    "email@veryhotmale.com",
		}

		user, err := service.UpdateOne(context.TODO(), updateParam)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.IsType(t, user, User{})
		assert.Equal(t, user.ID, idObj)
		assert.Equal(t, user.FullName, "John Doe")
		assert.Equal(t, user.Email, "email@veryhotmale.com")
		assert.Equal(t, user.CreatedAt, randomDate)
	})

	mt.Run("invalid object id returns empty user and err", func(mt *mtest.T) {
		service := NewService(mt.Coll)

		updateParam := UpdateParam{
			ID:       "123",
			FullName: "John Doe",
			Email:    "lmamo@h.com",
		}

		user, err := service.UpdateOne(context.TODO(), updateParam)

		assert.NotNil(t, err)
		assert.Empty(t, user)
		assert.IsType(t, user, User{})
	})

	mt.Run("encoding error returns empty user and err", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse())

		service := NewService(mt.Coll)

		updateParam := UpdateParam{
			ID:       "123",
			FullName: "John Doe",
			Email:    "la@la.com",
		}

		user, err := service.UpdateOne(context.TODO(), updateParam)

		assert.NotNil(t, err)
		assert.Empty(t, user)
		assert.IsType(t, user, User{})
	})
}
