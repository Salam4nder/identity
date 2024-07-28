package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Salam4nder/user/pkg/password"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("db")

// User defines a user in the users table.
type User struct {
	ID           uuid.UUID  `db:"id"`
	FullName     string     `db:"full_name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

// CreateUserParams defines the parameters to [CreateUser].
type CreateUserParams struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Password  password.SafeString
	CreatedAt time.Time
}

func (x CreateUserParams) SpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("user_id", x.ID.String()),
		attribute.String("full_name", x.FullName),
		attribute.String("email", x.Email),
		attribute.Int("password_length", len(x.Password)),
	}
}

// CreateUser creates a new user in the database.
func (x *SQL) CreateUser(ctx context.Context, params CreateUserParams) error {
	ctx, span := tracer.Start(ctx, "db.CreateUser", trace.WithAttributes(params.SpanAttributes()...))
	defer span.End()

	query := `
    INSERT INTO users (id, full_name, email, password_hash, created_at)
    VALUES ($1, $2, $3, $4, $5)
    `
	span.SetAttributes(attribute.String("query", query))

	res, err := x.db.ExecContext(
		ctx,
		query,
		params.ID,
		params.FullName,
		params.Email,
		params.Password,
		params.CreatedAt,
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		if IsSentinelErr(err) {
			return SentinelErr(err)
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
		switch rowsAffected {
		case 0:
			return ErrNoRowsAffected
		default:
			return ErrMultipleRowsAffected
		}
	}

	return nil
}

// ReadUser reads a user from the database.
func (x *SQL) ReadUser(ctx context.Context, id uuid.UUID) (*User, error) {
	ctx, span := tracer.Start(ctx, "db.ReadUser")
	defer span.End()
	span.SetAttributes(attribute.String("id", id.String()))

	if id == uuid.Nil {
		return nil, InputError{Field: "id", Value: id.String()}
	}

	query := `
        SELECT id, full_name, email, password_hash, created_at, updated_at
        FROM users
        WHERE id = $1
        `
	span.SetAttributes(attribute.String("query", query))

	var user User
	if err := x.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// ReadUserByEmail reads a user from the database by email.
func (x *SQL) ReadUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, span := tracer.Start(ctx, "db.ReadUserByEmail")
	defer span.End()

	if email == "" {
		return nil, InputError{Field: "email", Value: email}
	}

	query := `
        SELECT id, full_name, email, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1
        `
	span.SetAttributes(
		attribute.String("query", query),
		attribute.String("email", email),
	)

	var user User
	if err := x.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUserParams defines the parameters used to update a user.
type UpdateUserParams struct {
	ID       uuid.UUID
	FullName string
	Email    string
}

func (x UpdateUserParams) SpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("user_id", x.ID.String()),
		attribute.String("full_name", x.FullName),
		attribute.String("email", x.Email),
	}
}

// UpdateUser updates a user in the database.
func (x *SQL) UpdateUser(ctx context.Context, params UpdateUserParams) error {
	ctx, span := tracer.Start(ctx, "db.UpdateUser", trace.WithAttributes(params.SpanAttributes()...))
	defer span.End()

	query := `
        UPDATE users
        SET full_name = $1, email = $2, updated_at = $3
        WHERE id = $4
        `
	span.SetAttributes(attribute.String("query", query))

	res, err := x.db.ExecContext(
		ctx,
		query,
		params.FullName,
		params.Email,
		time.Now(),
		params.ID,
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		if IsSentinelErr(err) {
			return SentinelErr(err)
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
		switch rowsAffected {
		case 0:
			return ErrNoRowsAffected
		default:
			return ErrMultipleRowsAffected
		}
	}

	return nil
}

// DeleteUser deletes a user from the database.
func (x *SQL) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "db.DeleteUser")
	defer span.End()

	query := `
        DELETE FROM users
        WHERE id = $1
        `
	span.SetAttributes(
		attribute.String("user_id", id.String()),
		attribute.String("query", query),
	)

	res, err := x.db.ExecContext(ctx, query, id)
	if err != nil {
		if IsSentinelErr(err) {
			return SentinelErr(err)
		}
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		switch rowsAffected {
		case 0:
			return ErrNoRowsAffected
		default:
			return ErrMultipleRowsAffected
		}
	}

	return nil
}
