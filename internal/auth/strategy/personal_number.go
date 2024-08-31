package strategy

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"

	"github.com/Salam4nder/identity/internal/database/personalnumber"
	"github.com/Salam4nder/identity/proto/gen"
)

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

func (x *PersonalNumber) Register(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	n, err := generateRandomUint64()
	if err != nil {
		return err
	}

	if err := personalnumber.Insert(ctx, x.db, n); err != nil {
		return err
	}

	return nil
}

func (x *PersonalNumber) Authenticate(_ context.Context) error {
	return nil
}

func (x *PersonalNumber) Revoke(_ context.Context) error {
	return nil
}

func (x *PersonalNumber) Renew(_ context.Context) error {
	return nil
}

func generateRandomUint64() (uint64, error) {
	var result uint64
	for i := 0; i < 16; i++ {
		// Generate a random digit between 0 and 9.
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return 0, err
		}

		// Shift the result left by one digit and add the new digit.
		result = result*10 + digit.Uint64()
	}

	return result, nil
}
