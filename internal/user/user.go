package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FullName  string             `bson:"full_name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt string             `bson:"created_at"`
}

// Filter represents a filter for users.
type Filter struct {
	FullName  string `bson:"full_name, omitempty"`
	Email     string `bson:"email, omitempty"`
	CreatedAt string `bson:"created_at, omitempty"`
}

// CreateParam represents a parameter for creating a user.
type CreateParam struct {
	FullName  string `bson:"full_name"`
	Email     string `bson:"email"`
	Password  string `bson:"password"`
	CreatedAt string `bson:"created_at"`
}

// UpdateParam represents a parameter for updating a user.
type UpdateParam struct {
	ID       string `bson:"_id"`
	FullName string `bson:"full_name"`
	Email    string `bson:"email"`
}

// FindOneResponse represents a response for finding a user.
type FindOneResponse struct {
	ID        primitive.ObjectID `bson:"_id"`
	FullName  string             `bson:"full_name"`
	Email     string             `bson:"email"`
	CreatedAt string             `bson:"created_at"`
}
