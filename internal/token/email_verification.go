package token

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// NewEmailVerificationToken will return a token in the format of
// sha256SumOf(email)/uuid.New
func NewEmailVerificationToken(email string) string {
	b := sha256.Sum256([]byte(email))
	h := hex.EncodeToString(b[:])
	return fmt.Sprintf("%s/%s", h, email)
}
