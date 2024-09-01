package personalnumber

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"

	"github.com/Salam4nder/identity/internal/database/personalnumber"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	tracer = otel.Tracer("personalnumber")

	key ctxKey
)

// Strategy implements the [Strategy] interface and has everything
// to be able to [Register], [Authenticate] and [Revoke]
// with a personal number.
type (
	ctxKey int

	Strategy struct {
		db *sql.DB
	}
)

func New(db *sql.DB) *Strategy {
	return &Strategy{db: db}
}

func FromContext(ctx context.Context) (uint64, error) {
	c, ok := ctx.Value(key).(uint64)
	if !ok {
		return 0, errors.New("strategy: getting personal_number output from context")
	}
	return c, nil
}

func (x *Strategy) Register(ctx context.Context) (context.Context, error) {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	n, err := generateRandomUint64()
	if err != nil {
		return ctx, err
	}
	span.SetAttributes(attribute.Int64("generated_number", int64(n)))

	if err := personalnumber.Insert(ctx, x.db, n); err != nil {
		return ctx, err
	}

	return newContext(ctx, n), nil
}

func (x *Strategy) Authenticate(_ context.Context) error {
	return nil
}

func (x *Strategy) Revoke(_ context.Context) error {
	return nil
}

func (x *Strategy) Renew(_ context.Context) error {
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

func newContext(ctx context.Context, n uint64) context.Context {
	return context.WithValue(ctx, key, n)
}
