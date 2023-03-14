package user

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the database.
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

// InsertOneParam represents a parameter for creating a user.
type InsertOneParam struct {
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

// Validate validates the CreateParam parameter.
// Returns an error if the CreateParam parameter is invalid.
func (c *InsertOneParam) Validate() error {
	var (
		fullNameErr error
		emailErr    error
		passwordErr error
	)

	if c.FullName == "" {
		fullNameErr = errors.New("full name is required")
	}

	if c.Email == "" {
		emailErr = errors.New("email is required")
	}

	if c.Password == "" {
		passwordErr = errors.New("password is required")
	}

	return errors.Join(fullNameErr, emailErr, passwordErr)
}

// Validate validates the UpdateParam parameter.
// Returns an error if the UpdateParam parameter is invalid.
func (u *UpdateParam) Validate() error {
	if u.ID == "" {
		return errors.New("id is empty")
	}

	var (
		fullNameErr error
		emailErr    error
	)

	if u.FullName == "" {
		fullNameErr = errors.New("full name is empty")
	}

	if u.Email == "" {
		emailErr = errors.New("email is empty")
	}

	if fullNameErr != nil && emailErr != nil {
		return errors.Join(fullNameErr, emailErr)
	}

	return nil
}

// Validate validates the create parameter.
// Returns an error if the Filter parameteer is invalid.
func (f *Filter) Validate() error {
	if f.FullName == "" &&
		f.Email == "" &&
		f.CreatedAt == "" {
		return errors.New("filter is empty")
	}

	return nil
}
