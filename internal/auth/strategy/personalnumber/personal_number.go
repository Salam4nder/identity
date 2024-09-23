package personalnumber

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Salam4nder/identity/internal/database"
	"github.com/Salam4nder/identity/internal/database/personalnumber"
	"github.com/Salam4nder/identity/pkg/random"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	tracer = otel.Tracer("personalnumber")

	inputKey  ctxKey
	outputKey ctxKey

	ErrNumberNotFound = errors.New("personalnumber: number not found")
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

func NewContext(ctx context.Context, n uint64) context.Context {
	return context.WithValue(ctx, inputKey, n)
}

func FromContext(ctx context.Context) (uint64, error) {
	c, ok := ctx.Value(outputKey).(uint64)
	if !ok {
		return 0, errors.New("personalnumber: getting personal_number from context")
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

func (x *Strategy) Authenticate(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Authenticate")
	defer span.End()

	n, err := fromContext(ctx)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.Int64("number", int64(n)))

	_, err = personalnumber.Get(ctx, x.db, n)
	if err != nil {
		if errors.As(err, &database.NotFoundError{}) {
			return ErrNumberNotFound
		}
		return err
	}
	return nil
}

func newContext(ctx context.Context, n uint64) context.Context {
	return context.WithValue(ctx, outputKey, n)
}

func fromContext(ctx context.Context) (uint64, error) {
	v, ok := ctx.Value(inputKey).(uint64)
	if !ok {
		return 0, errors.New("personalnumber: getting number from context")
	}
	return v, nil
}
