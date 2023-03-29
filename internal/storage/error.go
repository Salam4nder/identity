package storage

import "errors"

// Common errors.
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidID      = errors.New("ID is not a valid ObjectID")
	ErrDuplicateEmail = errors.New("email already exists")
)
