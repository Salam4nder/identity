package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		min := int64(0)
		max := int64(100)
		n := RandomInt(min, max)

		assert.True(t, n >= min && n <= max)
	}
}

func TestRandomString(t *testing.T) {
	for i := 0; i < 100; i++ {
		length := 10
		s := RandomString(length)

		assert.Equal(t, length, len(s))
		assert.True(t, strings.ContainsAny(s, charset))
		assert.NotEmpty(t, s)
	}
}

func TestRandomEmail(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandomEmail()

		assert.True(t, strings.ContainsAny(s, charset))
		assert.True(t, strings.Contains(s, "@"))
		assert.True(t, strings.Contains(s, "."))
		assert.NotEmpty(t, s)
	}
}

func TestRandomFullName(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandomFullName()

		assert.True(t, strings.ContainsAny(s, charset))
		assert.True(t, strings.Contains(s, " "))
		assert.NotEmpty(t, s)
	}
}
