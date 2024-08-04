package strategy

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Salam4nder/identity/internal/auth"
	"github.com/Salam4nder/identity/internal/database/credentials"
	"github.com/Salam4nder/identity/internal/email"
	"github.com/Salam4nder/identity/pkg/password"
	"github.com/Salam4nder/identity/pkg/validation"
	"github.com/Salam4nder/identity/proto/gen"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("strategy")

var _ auth.Strategy = (*Credentials)(nil)

type (
	// Credentials implements the [Strategy] interface and has everything
	// to be able to [Register()], [Authenticate()] and [Revoke()] with credentials.
	Credentials struct {
		db       *sql.DB
		natsConn *nats.Conn

		email    string
		password password.SafeString
	}

	// CredentialsInput is the input for the credentials strategy.
	CredentialsInput struct {
		Email    string
		Password string
	}
)

func (x CredentialsInput) TraceAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("email", x.Email),
		attribute.Int("password length", utf8.RuneCountInString((x.Password))),
	}
}

// NewCredentials creates a new [Credentials] strategy for authentication.
// It still needs [CredentialsInput] to be able to call it's methods.
// An [CredentialsInput] is created with [IngestInput()].
func NewCredentials(db *sql.DB, natsConn *nats.Conn) *Credentials {
	return &Credentials{db: db, natsConn: natsConn}
}

func (x *Credentials) ConfiguredStrategy() gen.Strategy {
	return gen.Strategy_Credentials
}

// IngestInput sets the input field of the underlying [Credentials].
// [Credentials] is ready to call the rest of it's methods if this method returns no error.
func (x *Credentials) IngestInput(ctx context.Context, input CredentialsInput) error {
	_, span := tracer.Start(ctx, "IngestInput", trace.WithAttributes(input.TraceAttributes()...))
	defer span.End()

	p, err := password.FromString(input.Password)
	if err != nil {
		return fmt.Errorf("strategy: credentials, %w", err)
	}
	if err = validation.Email(input.Email); err != nil {
		return fmt.Errorf("strategy: credentials, %w", err)
	}

	x.email = input.Email
	x.password = p

	return nil
}

// Register will handles registration with the credentials strategy.
// It will insert a new [credentials.Entry] into the credentials table
// and send an email to the registered user.
func (x *Credentials) Register(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	if err := credentials.Insert(ctx, x.db, credentials.InsertParams{
		ID:        uuid.New(),
		Email:     x.email,
		Password:  x.password,
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}

	if err := email.Ingest(ctx, x.natsConn, email.Email{
		To:      x.email,
		From:    email.TestFrom,
		Subject: email.TestSubject,
		Body:    email.TestBody,
	}); err != nil {
		return err
	}

	return nil
}

func (x *Credentials) Authenticate(_ context.Context) error {
	return nil
}

func (x *Credentials) Revoke(_ context.Context) error {
	return nil
}

func (x *Credentials) Renew(_ context.Context) error {
	return nil
}
