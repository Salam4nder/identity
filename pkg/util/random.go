package util

import (
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

// RandomInt returns a random integer between min and max.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString returns a random string of length.
func RandomString(length int) string {
	var builder strings.Builder

	k := len(charset)

	for i := 0; i < length; i++ {
		c := charset[rand.Intn(k)]
		builder.WriteByte(c)
	}

	return builder.String()
}

// RandomEmail returns a random email address.
func RandomEmail() string {
	return RandomString(10) + "@" + RandomString(5) + ".com"
}

// RandomFullName returns a random full name.
func RandomFullName() string {
	return RandomString(10) + " " + RandomString(10)
}

// RandomDate returns a random date.
func RandomDate() time.Time {
	return time.Now().AddDate(0, 0, -int(RandomInt(0, 365)))
}
