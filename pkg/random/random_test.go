package random

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	for i := 0; i < 20; i++ {
		min := int64(0)
		max := int64(100)
		res := Int(min, max)

		assert.True(t, res >= min && res <= max)
	}
}

func TestString(t *testing.T) {
	for i := 0; i < 20; i++ {
		length := 10
		res := String(length)

		assert.Equal(t, length, len(res))
		assert.True(t, strings.ContainsAny(res, charset))
		assert.NotEmpty(t, res)
	}
}

func TestEmail(t *testing.T) {
	for i := 0; i < 20; i++ {
		res := Email()

		assert.True(t, strings.ContainsAny(res, charset))
		assert.True(t, strings.Contains(res, "@"))
		assert.True(t, strings.Contains(res, "."))
		assert.NotEmpty(t, res)
	}
}

func TestFullName(t *testing.T) {
	for i := 0; i < 20; i++ {
		res := FullName()

		assert.True(t, strings.ContainsAny(res, charset))
		assert.True(t, strings.Contains(res, " "))
		assert.NotEmpty(t, res)
	}
}

func TestDate(t *testing.T) {
	for i := 0; i < 20; i++ {
		res := Date()

		assert.NotEmpty(t, res)
		assert.IsType(t, res, time.Time{})
	}
}
