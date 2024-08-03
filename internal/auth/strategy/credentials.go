package strategy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Salam4nder/user/internal/auth"
	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/email"
	"github.com/Salam4nder/user/pkg/password"
	"github.com/Salam4nder/user/pkg/validation"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
)

var _ auth.Strategy = (*Credentials)(nil)

type (
	// Credentials implements the [Strategy] interface and has everything
	// to be able to [Register()], [Authenticate()] and [Revoke()] with credentials.
	Credentials struct {
		input    Input
		db       *sql.DB
		natsConn *nats.Conn
	}

	// Input is the input for the credentials strategy.
	Input struct {
		Email    string
		Password password.SafeString
	}
)

// NewCredentials creates a new [Credentials] strategy for authentication.
// It still needs [Input] to be able to call it's methods.
// An [Input] is created with [IngestInput()].
func NewCredentials(db *sql.DB, natsConn *nats.Conn) *Credentials {
	return &Credentials{db: db, natsConn: natsConn}
}

// IngestInput sets the input field of the underlying [Credentials].
// [Credentials] is ready to call the rest of it's methods if this method returns no error.
func (x *Credentials) IngestInput(ctx context.Context, email, pw string) error {
	ctx, span := tracer.Start(ctx, "IngestInput")
	defer span.End()
	span.SetAttributes(
		attribute.String("email", email),
		attribute.Int64("password length", int64(utf8.RuneCountInString(pw))),
	)

	p, err := password.FromString(pw)
	if err != nil {
		return fmt.Errorf("strategy: credentials password, %w", err)
	}
	if err = validation.Email(email); err != nil {
		return fmt.Errorf("strategy: credentials email, %w", err)
	}

	x.input = Input{
		Email:    email,
		Password: p,
	}

	return nil
}

func (x *Credentials) Register(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	if err := db.CreateUser(ctx, x.db, db.CreateUserParams{
		ID:        uuid.New(),
		FullName:  "Full Name",
		Email:     x.input.Email,
		Password:  x.input.Password,
		CreatedAt: time.Now(),
	}); err != nil {
		// TODO(kg:) Errs.
		if errors.Is(err, db.ErrDuplicateEmail) {
			return err
		}
		return err
	}

	if err := email.Ingest(ctx, x.natsConn, email.Email{
		To:      x.input.Email,
		From:    email.TestFrom,
		Subject: email.TestSubject,
		Body:    email.TestBody,
	}); err != nil {
		// TODO(kg:) Errs.
		return err
	}

	return nil
}
func (x *Credentials) Authenticate(ctx context.Context) error
func (x *Credentials) Revoke(ctx context.Context) error
