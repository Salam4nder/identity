package token

import "testing"

func TestSafeString(t *testing.T) {
	t.Run("stringer", func(t *testing.T) {
		if fromString("ass").String() != "REDACTED" {
			t.Error("expected redacted")
		}
	})
}
