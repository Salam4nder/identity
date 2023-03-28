package storage

// Error is used for custom errors.
type Error string

func (e Error) Error() string { return string(e) }

// Common errors
const (
	ErrUserNotFound = Error("user not found")
	ErrInvalidID    = Error("ID is not a valid ObjectID")
)
