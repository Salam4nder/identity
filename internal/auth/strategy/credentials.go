package strategy

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

	credentialsKey ctxKey
)

type (
	ctxKey int

	// Credentials implements the [Strategy] interface and has everything
	// to be able to [Register], [Authenticate] and [Revoke] with credentials.
	Credentials struct {
		db       *sql.DB
		natsConn *nats.Conn
	}
)

type Input struct {
	Email, Password string
}

// NewCredentials creates a new [Credentials] strategy for authentication.
// It still needs [CredentialsInput] to be able to call it's methods.
// A [CredentialsInput] is created with [IngestInput()].
func NewCredentials(db *sql.DB, natsConn *nats.Conn) *Credentials {
	return &Credentials{db: db, natsConn: natsConn}
}

func credentialsFromContext(ctx context.Context) (Input, error) {
	c, ok := ctx.Value(credentialsKey).(Input)
	if !ok {
		return Input{}, errors.New("strategy: getting credentials key from context")
	}
	return c, nil
}

func NewContext(ctx context.Context, c Input) context.Context {
	return context.WithValue(ctx, credentialsKey, c)
}

// Register will handles registration with the credentials strategy.
// It will insert a new [credentials.Entry] into the credentials table
// and send an email to the registered user.
func (x *Credentials) Register(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	cred, err := credentialsFromContext(ctx)
	if err != nil {
		return err
	}
	span.SetAttributes(
		attribute.String("email", cred.Email),
		attribute.Int("password length", utf8.RuneCountInString((cred.Password))),
	)

	p, err := password.FromString(cred.Password)
	if err != nil {
		return fmt.Errorf("strategy: credentials, %w", err)
	}
	if err = validation.Email(cred.Email); err != nil {
		return fmt.Errorf("strategy: credentials, %w", err)
	}

	if err := credentials.Insert(ctx, x.db, credentials.InsertParams{
		ID:        uuid.New(),
		Email:     cred.Email,
		Password:  p,
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}

	if err := email.Ingest(ctx, x.natsConn, email.Email{
		To:      cred.Email,
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
