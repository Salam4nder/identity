package credentials

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Salam4nder/user/internal/database"
	"github.com/Salam4nder/user/pkg/password"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("credentials")

const Tablename = "credentials"

// Entry defines an entry in the credentials table.
type Entry struct {
	ID           uuid.UUID  `db:"id"`
	FullName     string     `db:"full_name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

// InsertParams defines the parameters for inserts.
type InsertParams struct {
	ID        uuid.UUID
	Email     string
	Password  password.SafeString
	CreatedAt time.Time
}

func (x InsertParams) SpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("user_id", x.ID.String()),
		attribute.String("email", x.Email),
		attribute.Int("password_length", len(x.Password)),
	}
}

// Insert a new credentials entry. Returns [database.DuplicateEntryError] on duplicate entry.
func Insert(ctx context.Context, db *sql.DB, params InsertParams) error {
	ctx, span := tracer.Start(ctx, "Insert", trace.WithAttributes(params.SpanAttributes()...))
	defer span.End()

	query := `
    INSERT INTO credentials (id, email, password_hash, created_at)
    VALUES ($1, $2, $3, $4)
    `
	span.SetAttributes(attribute.String("query", query))

	res, err := db.ExecContext(
		ctx,
		query,
		params.ID,
		params.Email,
		params.Password,
		params.CreatedAt,
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		if database.IsPSQLDuplicateEntryError(err) {
			return database.NewDuplicateEntryError("credentials")
		}
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	if rowsAffected != 1 {
		return database.NewRowsAffectedError(1, rowsAffected)
	}

	return nil
}

// Read a credentials [Entry] by ID.
func Read(ctx context.Context, db *sql.DB, id uuid.UUID) (*Entry, error) {
	ctx, span := tracer.Start(ctx, "Read")
	defer span.End()
	span.SetAttributes(attribute.String("id", id.String()))

	if id == uuid.Nil {
		return nil, database.NewInputError("id", id.String())
	}

	query := `
        SELECT id, email, password_hash, created_at, updated_at
        FROM credentials
        WHERE id = $1
        `
	span.SetAttributes(attribute.String("query", query))

	var user Entry
	if err := db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.NewNotFoundError("credentials", id.String())
		}
		return nil, err
	}

	return &user, nil
}

// ReadByEmail a credentials [Entry] by an email.
func ReadByEmail(ctx context.Context, db *sql.DB, email string) (*Entry, error) {
	ctx, span := tracer.Start(ctx, "ReadByEmail")
	defer span.End()

	if email == "" {
		return nil, database.NewInputError("email", email)
	}

	query := `
        SELECT id, email, password_hash, created_at, updated_at
        FROM credentials
        WHERE email = $1
        `
	span.SetAttributes(
		attribute.String("query", query),
		attribute.String("email", email),
	)

	var user Entry
	if err := db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.NewNotFoundError("credentials", email)
		}
		return nil, err
	}

	return &user, nil
}

// UpdateParams defines the parameters used to update credentials.
type UpdateParams struct {
	ID    uuid.UUID
	Email string
}

func (x UpdateParams) SpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("user_id", x.ID.String()),
		attribute.String("email", x.Email),
	}
}

// Update credentials. Returns [database.DuplicateEntryError] on duplicate entry.
func Update(ctx context.Context, db *sql.DB, params UpdateParams) error {
	ctx, span := tracer.Start(ctx, "Update", trace.WithAttributes(params.SpanAttributes()...))
	defer span.End()

	query := `
        UPDATE credentials
        SET email = $1, updated_at = $2
        WHERE id = $3
        `
	span.SetAttributes(attribute.String("query", query))

	res, err := db.ExecContext(
		ctx,
		query,
		params.Email,
		time.Now(),
		params.ID,
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		if database.IsPSQLDuplicateEntryError(err) {
			return database.NewDuplicateEntryError("credentials")
		}
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	if rowsAffected != 1 {
		return database.NewRowsAffectedError(1, rowsAffected)
	}

	return nil
}

// Delete a credentils [Entry] from the database.
func Delete(ctx context.Context, db *sql.DB, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()

	query := `
        DELETE FROM credentials
        WHERE id = $1
        `
	span.SetAttributes(
		attribute.String("user_id", id.String()),
		attribute.String("query", query),
	)

	res, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return database.NewRowsAffectedError(1, rowsAffected)
	}

	return nil
}
