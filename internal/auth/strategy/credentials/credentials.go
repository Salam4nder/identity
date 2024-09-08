package credentials

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Salam4nder/identity/internal/database/credentials"
	"github.com/Salam4nder/identity/internal/email"
	"github.com/Salam4nder/identity/pkg/password"
	"github.com/Salam4nder/identity/pkg/validation"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	tracer = otel.Tracer("strategy")

	inputKey  ctxKey
	outputKey ctxKey
)

type (
	ctxKey int

	// Strategy implements the [Strategy] interface and has everything
	// to be able to [Register] and [Authenticate] with credentials.
	Strategy struct {
		db       *sql.DB
		natsConn *nats.Conn
	}

	Input struct {
		Email, Password string
	}

	Output struct {
		Email string
	}
)

// New creates a new [Strategy] for authentication.
func New(db *sql.DB, natsConn *nats.Conn) *Strategy {
	return &Strategy{db: db, natsConn: natsConn}
}

func NewContext(ctx context.Context, c *Input) context.Context {
	return context.WithValue(ctx, inputKey, c)
}

func FromContext(ctx context.Context) (*Output, error) {
	c, ok := ctx.Value(outputKey).(*Output)
	if !ok {
		return nil, errors.New("strategy: getting credentials output from context")
	}
	return c, nil
}

// Register will handles registration with the credentials strategy.
// It will insert a new [credentials.Entry] into the credentials table
// and send an email to the registered user.
func (x *Strategy) Register(ctx context.Context) (context.Context, error) {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	cred, err := fromContext(ctx)
	if err != nil {
		return ctx, err
	}
	span.SetAttributes(
		attribute.String("email", cred.Email),
		attribute.Int("password length", utf8.RuneCountInString((cred.Password))),
	)

	p, err := password.FromString(cred.Password)
	if err != nil {
		return ctx, fmt.Errorf("strategy: credentials, %w", err)
	}
	if err = validation.Email(cred.Email); err != nil {
		return ctx, fmt.Errorf("strategy: credentials, %w", err)
	}

	if err := credentials.Insert(ctx, x.db, credentials.InsertParams{
		ID:        uuid.New(),
		Email:     cred.Email,
		Password:  p,
		CreatedAt: time.Now(),
	}); err != nil {
		return ctx, err
	}

	if err := email.Ingest(ctx, x.natsConn, email.Email{
		To:      cred.Email,
		From:    email.TestFrom,
		Subject: email.TestSubject,
		Body:    email.TestBody,
	}); err != nil {
		return ctx, err
	}

	return newContext(ctx, &Output{Email: cred.Email}), nil
}

func (x *Strategy) Authenticate(_ context.Context) error {
	return nil
}

func fromContext(ctx context.Context) (*Input, error) {
	c, ok := ctx.Value(inputKey).(*Input)
	if !ok {
		return nil, errors.New("strategy: getting credentials input from context")
	}
	return c, nil
}

func newContext(ctx context.Context, c *Output) context.Context {
	return context.WithValue(ctx, outputKey, c)
}
