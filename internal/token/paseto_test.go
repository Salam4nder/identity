package token

import (
	"testing"
	"time"

	"github.com/Salam4nder/identity/pkg/random"
	"github.com/Salam4nder/identity/proto/gen"
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
	s, err := b.MakeAccessToken(random.Email(), gen.Strategy_TypeCredentials)
	if err != nil {
		t.Error("expected no error")
	}
	if s == "" {
		t.Error("token is empty")
	}

	t.Run("bad identifier", func(t *testing.T) {
		b := bootstrap(t)
		s, err := b.MakeAccessToken(100, gen.Strategy_TypeCredentials)
		if err == nil {
			t.Error("expected error")
		}
		if s != "" {
			t.Error("expected empty string")
		}
	})
	t.Run("bad strategy", func(t *testing.T) {
		b := bootstrap(t)
		s, err := b.MakeAccessToken(100, gen.Strategy_TypeNoStrategy)
		if err == nil {
			t.Error("expected error")
		}
		if s != "" {
			t.Error("expected empty string")
		}
	})
}

func TestMakeRefreshToken(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		b := bootstrap(t)
		s, err := b.MakeRefreshToken(random.Email(), gen.Strategy_TypeCredentials)
		if err != nil {
			t.Error("expected no error")
		}
		if s == "" {
			t.Error("token is empty")
		}
	})
	t.Run("bad identifier", func(t *testing.T) {
		b := bootstrap(t)
		s, err := b.MakeRefreshToken(100, gen.Strategy_TypeCredentials)
		if err == nil {
			t.Error("expected error")
		}
		if s != "" {
			t.Error("expected empty string")
		}
	})
	t.Run("bad strategy", func(t *testing.T) {
		b := bootstrap(t)
		s, err := b.MakeRefreshToken(100, gen.Strategy_TypeNoStrategy)
		if err == nil {
			t.Error("expected error")
		}
		if s != "" {
			t.Error("expected empty string")
		}
	})
}

func TestVerify(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		b := bootstrap(t)
		s, err := b.MakeAccessToken(random.Email(), gen.Strategy_TypeCredentials)
		if err != nil {
			t.Error("expected no error")
		}
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
