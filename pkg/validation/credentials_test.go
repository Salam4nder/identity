package validation

import (
	"errors"
	"strings"
	"testing"
)

func TestStringLen(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		if err := StringLen(strings.Repeat("a", 100), 1, 100, "test"); err != nil {
			t.Errorf("ValidateString failed: %s", err)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		if err := StringLen("", 1, 100, "test"); err == nil {
			t.Errorf("ValidateString failed: %s", err)
		}
	})

	t.Run("Too short", func(t *testing.T) {
		if err := StringLen("a", 2, 100, "test"); err == nil {
			t.Errorf("ValidateString(\"a\") failed: %s", err)
		}
	})
}

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
		if !errors.Is(err, ErrInvalidUsername) {
			t.Errorf("ValidateUsername failed: %s", err)
		}
	})
}

func TestPassword(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		if err := Password(""); err == nil {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Too short", func(t *testing.T) {
		if err := Password("a"); err == nil {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Too long", func(t *testing.T) {
		if err := Password(strings.Repeat("a", 101)); err == nil {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Does not contain uppercase", func(t *testing.T) {
		err := Password("password1")
		if !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Does not contain lowercase", func(t *testing.T) {
		err := Password("PASSWORD1")
		if !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("Does not contain number", func(t *testing.T) {
		err := Password("Password")
		if !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("ValidatePassword failed: %s", err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		err := Password("Password1")
		if err != nil {
			t.Errorf("ValidatePassword failed: %s", err)
		}
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
		if !errors.Is(err, ErrInvalidEmail) {
			t.Errorf("ValidateEmail failed: %s", err)
		}
	})
}

func TestSecret(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		if err := Secret(strings.Repeat("a", 33)); err != nil {
			t.Errorf("ValidateSecret failed: %s", err)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		if err := Secret(""); err == nil {
			t.Errorf("ValidateSecret failed: %s", err)
		}
	})

	t.Run("Too short", func(t *testing.T) {
		if err := Secret("a"); err == nil {
			t.Errorf("ValidateSecret failed: %s", err)
		}
	})
}
