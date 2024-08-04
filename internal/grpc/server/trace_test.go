package server

import (
	"testing"

	"github.com/Salam4nder/user/proto/gen"
)

func TestGenSpanAttributes(t *testing.T) {
	t.Run("gen.Credentials", func(t *testing.T) {
		c := &gen.Credentials{
			Email:    "email",
			Password: "password",
		}
		attr, err := GenSpanAttributes(c)
		if err != nil {
			t.Error("expected no error")
		}
		if len(attr) != 2 {
			t.Errorf("expected len 2, got %d", len(attr))
		}
	})

	t.Run("gen.Number", func(t *testing.T) {
		c := &gen.Number{
			Numbers: 23409820394802,
		}
		attr, err := GenSpanAttributes(c)
		if err != nil {
			t.Error("expected no error")
		}
		if len(attr) != 1 {
			t.Errorf("expected len 2, got %d", len(attr))
		}
	})
}
