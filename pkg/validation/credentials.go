package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString

	// ErrInvalidEmail is returned when the email is invalid.
	ErrInvalidEmail = errors.New("validation: provided email is not a valid email address")
	// ErrInvalidUsername is returned when the username is invalid.
	ErrInvalidUsername = errors.New("validation: username must contain only letters, digits, or underscore")
	// ErrInvalidFullName is returned when the full name is invalid.
	ErrInvalidFullName = errors.New("validation: full name must contain only letters or spaces")
	// ErrInvalidPassword is returned when the password is invalid.
	ErrInvalidPassword = errors.New("validation: password must contain an uppercase and lowercase letter and a digit")
)

// StringLen checks if the given string is between the given min and max.
func StringLen(value string, minLength int, maxLength int, entity string) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf(
			"validation: %s must contain from %d-%d characters",
			entity,
			minLength,
			maxLength,
		)
	}

	return nil
}

// Username checks if the given username is valid.
func Username(value string) error {
	if err := StringLen(value, 3, 100, "username"); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return ErrInvalidUsername
	}

	return nil
}

// FullName checks if the given full name is valid.
func FullName(value string) error {
	if err := StringLen(value, 3, 100, "full_name"); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return ErrInvalidFullName
	}

	return nil
}

// Password checks if the given password is valid.
func Password(value string) error {
	if err := StringLen(value, 6, 100, "password"); err != nil {
		return err
	}

	// Password must contain at least one uppercase letter,
	// one lowercase letter and one digit.
	if !regexp.MustCompile(`[a-z]`).MatchString(value) ||
		!regexp.MustCompile(`[A-Z]`).MatchString(value) ||
		!regexp.MustCompile(`[0-9]`).MatchString(value) {
		return ErrInvalidPassword
	}

	return nil
}

// Email checks if the given email is valid.
func Email(value string) error {
	if err := StringLen(value, 3, 200, "email"); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return ErrInvalidEmail
	}

	return nil
}

// Secret checks if the given secret is valid.
func Secret(value string) error {
	return StringLen(value, 32, 128, "secret")
}
