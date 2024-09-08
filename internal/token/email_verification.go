package token

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// NewEmailVerificationToken will return a token in the format of id/uuid.New.
func NewEmailVerificationToken(id string) string {
	return fmt.Sprintf("%s/%s", id, uuid.NewString())
}

func ParseVerificationToken(t string) (id uuid.UUID, token uuid.UUID, err error) {
	s := strings.Split(t, "/")
	if len(s) != 2 {
		return uuid.Nil, uuid.Nil, errors.New("token: len of token is not 2")
	}

	id, err = uuid.Parse(s[0])
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("token: parsing id, %w", err)
	}

	token, err = uuid.Parse(s[1])
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("token: parsing token, %w", err)
	}

	return id, token, err
}
