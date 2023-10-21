package util

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateStringLen(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		if err := ValidateStringLen(strings.Repeat("a", 100), 1, 100); err != nil {
			t.Errorf("ValidateString failed: %s", err)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		if err := ValidateStringLen("", 1, 100); err == nil {
			t.Errorf("ValidateString failed: %s", err)
		}
	})

	t.Run("Too short", func(t *testing.T) {
		if err := ValidateStringLen("a", 2, 100); err == nil {
			t.Errorf("ValidateString(\"a\") failed: %s", err)
		}
	})
}

func TestValidateUsername(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		if err := ValidateUsername("userName"); err != nil {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		if err := ValidateUsername(""); err == nil {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})

	t.Run("Too short", func(t *testing.T) {
		if err := ValidateUsername("a"); err == nil {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})

	t.Run("Too long", func(t *testing.T) {
		if err := ValidateUsername(strings.Repeat("a", 101)); err == nil {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})

	t.Run("Invalid character", func(t *testing.T) {
		err := ValidateUsername("username!")
		if !errors.Is(err, ErrInvalidUsername) {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		if err := ValidatePassword(""); err == nil {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Too short", func(t *testing.T) {
		if err := ValidatePassword("a"); err == nil {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Too long", func(t *testing.T) {
		if err := ValidatePassword(strings.Repeat("a", 101)); err == nil {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Does not contain uppercase", func(t *testing.T) {
		err := ValidatePassword("password1")
		if !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Does not contain lowercase", func(t *testing.T) {
		err := ValidatePassword("PASSWORD1")
		if !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Does not contain number", func(t *testing.T) {
		err := ValidatePassword("Password")
		if !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		err := ValidatePassword("Password1")
		if err != nil {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})
}
