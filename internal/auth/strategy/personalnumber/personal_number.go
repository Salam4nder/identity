package personalnumber

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Salam4nder/identity/internal/database/personalnumber"
	"github.com/Salam4nder/identity/pkg/random"
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

	n, err := random.UINT64()
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

func newContext(ctx context.Context, n uint64) context.Context {
	return context.WithValue(ctx, key, n)
}
