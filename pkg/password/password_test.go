package password

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestFromString(t *testing.T) {
	t.Run("String is longer than 73 bytes)", func(t *testing.T) {
		tooLong := strings.Repeat("a", 73)
		if len(tooLong) != 73 {
			t.Error("string is not 73 bytes")
		}
		_, err := FromString(tooLong)
		if !errors.As(err, &TooLongError{}) {
			t.Errorf("expected TooLongError, got %T", err)
		}
	})

	t.Run("String is shorter than 8 chars)", func(t *testing.T) {
		tooShort := strings.Repeat("a", 7)
		if utf8.RuneCountInString(tooShort) != 7 {
			t.Error("string is not 7 utf8 runes")
		}
		_, err := FromString(tooShort)
		if !errors.As(err, &TooShortError{}) {
			t.Errorf("expected TooShortError, got %T", err)
		}
	})

	t.Run("String is does not contain upper, lower and digit)", func(t *testing.T) {
		_, err := FromString("myAssWord")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("OK", func(t *testing.T) {
		_, err := FromString("myC00lp4zzW0rd")
		if err != nil {
			t.Error("expected no error")
		}
	})
}

func TestString(t *testing.T) {
	pw, err := FromString("myC00lp4zzW0rd")
	if err != nil {
		t.Error("expected no error")
	}

	if pw.String() != "REDACTED" {
		t.Errorf("expected REDACTED, got %s", pw.String())
	}
}
