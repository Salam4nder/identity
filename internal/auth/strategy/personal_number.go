package strategy

import (
	"database/sql"

	"github.com/Salam4nder/identity/internal/auth"
	"github.com/Salam4nder/identity/proto/gen"
)

var _ auth.Strategy = (*Credentials)(nil)

// PersonalNumber implements the [Strategy] interface and has everything
// to be able to [Register()], [Authenticate()] and [Revoke()]
// with a personal number.
type PersonalNumber struct {
	db *sql.DB

	number uint64
}

func NewPersonalNumber(db *sql.DB) *PersonalNumber {
	return &PersonalNumber{db: db}
}

func (x *PersonalNumber) ConfiguredStrategy() gen.Strategy {
	return gen.Strategy_PersonalNumber
}
