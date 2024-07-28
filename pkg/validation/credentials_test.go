package validation

import (
	"errors"
	"strings"
	"testing"
)

func TestUsername(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		if err := Username("userName"); err != nil {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		if err := Username(""); err == nil {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})

	t.Run("Too short", func(t *testing.T) {
		if err := Username("a"); err == nil {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})

	t.Run("Too long", func(t *testing.T) {
		if err := Username(strings.Repeat("a", 101)); err == nil {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})

	t.Run("Invalid character", func(t *testing.T) {
		err := Username("username!")
		if !errors.As(err, &InputError{}) {
			t.Error("expected InputError", err)
		}
	})
}

func TestFullName(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		t.Run("With space", func(t *testing.T) {
			err := FullName("John Adams")
			if err != nil {
				t.Errorf("Expected no error")
			}
		})
		t.Run("With dash", func(t *testing.T) {
			err := FullName("John Von-Haartman")
			if err != nil {
				t.Errorf("Expected no error")
			}
		})
	})
}

func TestEmail(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		if err := Email("email@email.com"); err != nil {
			t.Errorf("ValidateEmail failed: %s", err)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		if err := Email(""); err == nil {
			t.Errorf("ValidateEmail failed: %s", err)
		}
	})

	t.Run("Too short", func(t *testing.T) {
		if err := Email("a"); err == nil {
			t.Errorf("ValidateEmail failed: %s", err)
		}
	})

	t.Run("Too long", func(t *testing.T) {
		if err := Email(strings.Repeat("a", 300)); err == nil {
			t.Errorf("ValidateEmail failed: %s", err)
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		err := Email("email")
		if !errors.As(err, &InputError{}) {
			t.Error("expected InputError", err)
		}
	})
}
