package server

import (
	"testing"
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
}
