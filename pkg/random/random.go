package random

import (
	crypto "crypto/rand"
	"math/big"
	"math/rand/v2"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

// Int returns a random integer between min and max.
func Int(min, max int64) int64 {
	// nolint:gosec
	return min + rand.Int64N(max-min+1)
}

// String returns a random string of length.
func String(length int) string {
	var builder strings.Builder

	k := len(charset)

	for i := 0; i < length; i++ {
		// nolint:gosec
		c := charset[rand.IntN(k)]
		builder.WriteByte(c)
	}

	return builder.String()
}

// Email returns a random email address.
func Email() string {
	return String(10) + "@" + String(5) + ".com"
}

// FullName returns a random full name.
func FullName() string {
	return String(10) + " " + String(10)
}

// Date returns a random date.
func Date() time.Time {
	return time.Now().AddDate(0, 0, -int(Int(0, 365)))
}

func UINT64() (uint64, error) {
	var result uint64

	for range 16 {
		// Generate a random digit between 0 and 9.
		digit, err := crypto.Int(crypto.Reader, big.NewInt(10))
		if err != nil {
			return 0, err
		}

		// Shift the result left by one digit and add the new digit.
		result = result*10 + digit.Uint64()
	}

	return result, nil
}
