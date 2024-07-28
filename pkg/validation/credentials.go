package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode/utf8"
)

const (
	MaxNameLen = 100

	MinFullNameLen = 3
	MinUserNameLen = 2
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s\-]+$`).MatchString
)

type InputError struct {
	text string
}

func (x InputError) Error() string {
	return x.text
}

// Username checks if the given username is valid.
func Username(value string) error {
	n := utf8.RuneCountInString(value)
	if n > MaxNameLen || n < MinUserNameLen {
		return InputError{
			text: fmt.Sprintf(
				"validation: username must be between %d and %d characters",
				MinUserNameLen,
				MaxNameLen,
			),
		}
	}

	if !isValidUsername(value) {
		return InputError{text: "validation: username must contain only letters, digits, or underscore"}
	}

	return nil
}

// FullName checks if the given full name is valid.
func FullName(value string) error {
	n := utf8.RuneCountInString(value)
	if n > MaxNameLen || n < MinFullNameLen {
		return InputError{
			text: fmt.Sprintf(
				"validation: full name must be between %d and %d characters",
				MinFullNameLen,
				MaxNameLen,
			),
		}
	}

	if !isValidFullName(value) {
		return InputError{text: "validation: full name must contain only letters, dashes or spaces"}
	}

	return nil
}

// Email checks if the given email is valid.
func Email(value string) error {
	if _, err := mail.ParseAddress(value); err != nil {
		return InputError{
			text: "validation: provided email is not a valid email address",
		}
	}

	return nil
}
