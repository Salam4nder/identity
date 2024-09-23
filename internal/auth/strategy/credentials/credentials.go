package credentials

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"unicode/utf8"

	"github.com/Salam4nder/identity/internal/database"
	"github.com/Salam4nder/identity/internal/database/credentials"
	tokendb "github.com/Salam4nder/identity/internal/database/token"
	"github.com/Salam4nder/identity/internal/email"
	"github.com/Salam4nder/identity/internal/token"
	"github.com/Salam4nder/identity/pkg/password"
	"github.com/Salam4nder/identity/pkg/validation"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

var (
	tracer = otel.Tracer("strategy")

	inputKey  ctxKey
	outputKey ctxKey

	ErrUserNotFound      = errors.New("credentials: user does not exist")
	ErrUserNotVerified   = errors.New("credentials: user is not verified")
	ErrTokenDoesNotExist = errors.New("credentials: token does not exist")
	ErrIncorrectPassword = errors.New("credentials: incorrect password")
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
		return nil, errors.New("credentials: getting credentials output from context")
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
		return ctx, fmt.Errorf("credentials: creating password from string, %w", err)
	}
	if err = validation.Email(cred.Email); err != nil {
		return ctx, fmt.Errorf("credentials: validating email, %w", err)
	}

	tx, err := x.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			slog.ErrorContext(ctx, "credentials: error occurred, rolling back", "err", err)
			if err := tx.Rollback(); err != nil {
				slog.ErrorContext(ctx, "credentials: failed rollback", "err", err)
			}
		}
	}()

	id := uuid.New()
	if err = credentials.Insert(ctx, tx, credentials.InsertParams{
		ID:        id,
		Email:     cred.Email,
		Password:  p,
		CreatedAt: time.Now(),
	}); err != nil {
		return ctx, err
	}

	t := token.NewEmailVerificationToken(id.String())
	if err = tokendb.Insert(ctx, tx, t); err != nil {
		return ctx, err
	}

	if err = email.Ingest(ctx, x.natsConn, email.Email{
		To:      cred.Email,
		From:    email.TestFrom,
		Subject: email.TestSubject,
		Body:    email.Verification("https://example.com", t),
	}); err != nil {
		return ctx, err
	}

	if err = tx.Commit(); err != nil {
		slog.ErrorContext(ctx, "credentials: commit failed", "err", err)
		return ctx, err
	}

	return newContext(ctx, &Output{Email: cred.Email}), nil
}

// Authenticate will authenticate a user.
// Possible errors are [ErrIncorrectPassword], [ErrUserNotFound], [ErrUserNotVerified] and a wrapped error
// indicating an internal error.
func (x *Strategy) Authenticate(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Authenticate")
	defer span.End()

	cred, err := fromContext(ctx)
	if err != nil {
		return err
	}
	span.SetAttributes(
		attribute.String("email", cred.Email),
		attribute.Int("password length", utf8.RuneCountInString((cred.Password))),
	)

	p, err := password.FromString(cred.Password)
	if err != nil {
		return ErrIncorrectPassword
	}

	e, err := credentials.ReadByEmail(ctx, x.db, cred.Email)
	if err != nil {
		if errors.As(err, &database.NotFoundError{}) {
			return ErrUserNotFound
		}
		return fmt.Errorf("credentials: reading by email, %w", err)
	}

	if e.VerifiedAt == nil || e.VerifiedAt.IsZero() {
		return ErrUserNotVerified
	}

	if err = bcrypt.CompareHashAndPassword([]byte(e.PasswordHash), []byte(p)); err != nil {
		return fmt.Errorf("credentials: comparing password hash, %w", err)
	}

	return nil
}

func (x *Strategy) VerifyEmail(ctx context.Context, tokenInput string) error {
	ctx, span := tracer.Start(ctx, "VerifyEmail")
	defer span.End()

	t, err := tokendb.Get(ctx, x.db, tokenInput)
	if err != nil {
		if errors.As(err, &database.NotFoundError{}) {
			return ErrTokenDoesNotExist
		}
		return fmt.Errorf("credentials: getting token, %w", err)
	}

	id, _, err := token.ParseVerificationToken(t)
	if err != nil {
		return fmt.Errorf("credentials: pasring verification token, %w", err)
	}

	c, err := credentials.Read(ctx, x.db, id)
	if err != nil {
		return fmt.Errorf("credentials: reading user by email, %w", err)
	}
	if c.VerifiedAt != nil && !c.VerifiedAt.IsZero() {
		return errors.New("credentials: user already verified")
	}

	tx, err := x.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			slog.ErrorContext(ctx, "credentials: error occurred, rolling back", "err", err)
			if err := tx.Rollback(); err != nil {
				slog.ErrorContext(ctx, "credentials: failed rollback", "err", err)
			}
		}
	}()

	if err = credentials.Verify(ctx, tx, c.ID); err != nil {
		return fmt.Errorf("credentials: verifying user, %w", err)
	}

	if err = tokendb.Delete(ctx, tx, t); err != nil {
		return fmt.Errorf("credentials: deleting token, %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("credentials: committing transaction, %w", err)
	}

	return nil
}

func fromContext(ctx context.Context) (*Input, error) {
	c, ok := ctx.Value(inputKey).(*Input)
	if !ok {
		return nil, errors.New("credentials: getting credentials input from context")
	}
	return c, nil
}

func newContext(ctx context.Context, c *Output) context.Context {
	return context.WithValue(ctx, outputKey, c)
}
