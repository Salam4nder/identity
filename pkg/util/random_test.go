package util

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	for i := 0; i < 20; i++ {
		min := int64(0)
		max := int64(100)
		res := RandomInt(min, max)

		assert.True(t, res >= min && res <= max)
	}
}

func TestRandomString(t *testing.T) {
	for i := 0; i < 20; i++ {
		length := 10
		res := RandomString(length)

		assert.Equal(t, length, len(res))
		assert.True(t, strings.ContainsAny(res, charset))
		assert.NotEmpty(t, res)
	}
}

func TestRandomEmail(t *testing.T) {
	for i := 0; i < 20; i++ {
		res := RandomEmail()

		assert.True(t, strings.ContainsAny(res, charset))
		assert.True(t, strings.Contains(res, "@"))
		assert.True(t, strings.Contains(res, "."))
		assert.NotEmpty(t, res)
	}
}

func TestRandomFullName(t *testing.T) {
	for i := 0; i < 20; i++ {
		res := RandomFullName()

		assert.True(t, strings.ContainsAny(res, charset))
		assert.True(t, strings.Contains(res, " "))
		assert.NotEmpty(t, res)
	}
}

func TestRandomDate(t *testing.T) {
	for i := 0; i < 20; i++ {
		res := RandomDate()

		assert.NotEmpty(t, res)
		assert.IsType(t, res, time.Time{})
	}
}
