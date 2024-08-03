package server

import (
	"testing"

	"github.com/Salam4nder/user/pkg/random"
	"github.com/Salam4nder/user/proto/gen"
)

func TestGenSpanAttributes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		attrs, err := GenSpanAttributes(nil)
		if err == nil {
			t.Log("GenSpanAttributes() error = nil")
			t.Fail()
			return
		}
		if len(attrs) != 0 {
			t.Errorf("GenSpanAttributes() got = %v, want = 0", len(attrs))
		}
	})
	t.Run("UserID", func(t *testing.T) {
		req := &gen.UserID{
			Id: "1",
		}
		attrs, err := GenSpanAttributes(req)
		if err != nil {
			t.Errorf("GenSpanAttributes() error = %v", err)
			return
		}
		if len(attrs) != 1 {
			t.Errorf("GenSpanAttributes() got = %v, want = 1", len(attrs))
		}
	})
	t.Run("CreateUserRequest", func(t *testing.T) {
		req := &gen.CreateUserRequest{
			FullName: random.FullName(),
			Email:    random.Email(),
			Password: random.String(8),
		}
		attrs, err := GenSpanAttributes(req)
		if err != nil {
			t.Errorf("GenSpanAttributes() error = %v", err)
			return
		}
		if len(attrs) != 3 {
			t.Errorf("GenSpanAttributes() got = %v, want = 3", len(attrs))
		}
	})
	t.Run("UpdateUserRequest", func(t *testing.T) {
		req := &gen.UpdateUserRequest{
			Id:       "1",
			FullName: random.FullName(),
			Email:    random.Email(),
		}
		attrs, err := GenSpanAttributes(req)
		if err != nil {
			t.Errorf("GenSpanAttributes() error = %v", err)
			return
		}
		if len(attrs) != 3 {
			t.Errorf("GenSpanAttributes() got = %v, want = 3", len(attrs))
		}
	})
	t.Run("ReadByEmailRequest", func(t *testing.T) {
		req := &gen.ReadByEmailRequest{
			Email: random.Email(),
		}
		attrs, err := GenSpanAttributes(req)
		if err != nil {
			t.Errorf("GenSpanAttributes() error = %v", err)
			return
		}
		if len(attrs) != 1 {
			t.Errorf("GenSpanAttributes() got = %v, want = 1", len(attrs))
		}
	})
	t.Run("LoginUserRequest", func(t *testing.T) {
		req := &gen.LoginUserRequest{
			Email:    random.Email(),
			Password: random.String(8),
		}
		attrs, err := GenSpanAttributes(req)
		if err != nil {
			t.Errorf("GenSpanAttributes() error = %v", err)
			return
		}
		if len(attrs) != 2 {
			t.Errorf("GenSpanAttributes() got = %v, want = 2", len(attrs))
		}
	})
}
