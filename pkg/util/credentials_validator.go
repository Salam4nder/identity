package util

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString

	// ErrInvalidEmail is returned when the email is invalid.
	ErrInvalidEmail = fmt.Errorf(
		"util: email validation failed, provided email is not a valid email address",
	)
	// ErrInvalidUsername is returned when the username is invalid.
	ErrInvalidUsername = fmt.Errorf(
		"util: username validaton failed, " +
			"must contain only letters, digits, or underscore",
	)
	// ErrInvalidFullName is returned when the full name is invalid.
	ErrInvalidFullName = fmt.Errorf(
		"util: full name validation failed, must contain only letters or spaces",
	)
	// ErrInvalidPassword is returned when the password is invalid.
	ErrInvalidPassword = fmt.Errorf(
		"util: password validation failed, " +
			"must contain at least one uppercase letter, one lowercase letter and one digit",
	)
)

// ValidateStringLen checks if the given string is between the given min and max.
func ValidateStringLen(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf(
			"util: validation failed, must contain from %d-%d characters",
			minLength,
			maxLength,
		)
	}

	return nil
}

// ValidateUsername checks if the given username is valid.
func ValidateUsername(value string) error {
	if err := ValidateStringLen(value, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return ErrInvalidUsername
	}

	return nil
}

// ValidateFullName checks if the given full name is valid.
func ValidateFullName(value string) error {
	if err := ValidateStringLen(value, 3, 100); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return ErrInvalidFullName
	}

	return nil
}

// ValidatePassword checks if the given password is valid.
func ValidatePassword(value string) error {
	if err := ValidateStringLen(value, 6, 100); err != nil {
		return err
	}

	// Password must contain at least one uppercase letter,
	// one lowercase letter and one digit
	if !regexp.MustCompile(`[a-z]`).MatchString(value) ||
		!regexp.MustCompile(`[A-Z]`).MatchString(value) ||
		!regexp.MustCompile(`[0-9]`).MatchString(value) {
		return ErrInvalidPassword
	}

	return nil
}

// ValidateEmail checks if the given email is valid.
func ValidateEmail(value string) error {
	if err := ValidateStringLen(value, 3, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return ErrInvalidEmail
	}

	return nil
}

// ValidateSecret checks if the given secret is valid.
func ValidateSecret(value string) error {
	return ValidateStringLen(value, 32, 128)
}
