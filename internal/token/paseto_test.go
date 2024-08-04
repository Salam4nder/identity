package token

import (
	"testing"
	"time"
)

func bootstrap(t *testing.T) *PasetoMaker {
	t.Helper()

	t.Run("invalid symmetric key", func(t *testing.T) {
		bb := make([]byte, 0, 31)
		for range 31 {
			bb = append(bb, byte('s'))
		}
		_, err := BootstrapPasetoMaker(
			time.Second*10,
			time.Minute,
			bb,
		)
		if err == nil {
			t.Error("expected err with invalid symmetric key")
		}
	})

	b := make([]byte, 0, 32)
	for range 32 {
		b = append(b, byte('s'))
	}
	maker, err := BootstrapPasetoMaker(
		time.Second*10,
		time.Minute,
		b,
	)
	if err != nil {
		t.Errorf("expected no err, got %s", err.Error())
	}
	return maker
}

func TestMakeAccessToken(t *testing.T) {
	b := bootstrap(t)
	s := b.MakeAccessToken()
	if s == "" {
		t.Error("token is empty")
	}
}

func TestMakeRefreshToken(t *testing.T) {
	b := bootstrap(t)
	s := b.MakeRefreshToken()
	if s == "" {
		t.Error("token is empty")
	}
}

func TestVerify(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		b := bootstrap(t)
		s := b.MakeAccessToken()
		if err := b.Verify((s)); err != nil {
			t.Errorf("expected no error, got %s", err.Error())
		}
	})

	t.Run("invalid returns error", func(t *testing.T) {
		b := bootstrap(t)
		if err := b.Verify(fromString("ass")); err == nil {
			t.Error("expected error")
		}
	})
}
