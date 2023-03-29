package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// ComparePasswordHash validates a password using bcrypt.
func ComparePasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(password))
}
