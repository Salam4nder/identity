package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

// User defines a user in the database.
type User struct {
	ID           uuid.UUID  `db:"id"`
	FullName     string     `db:"full_name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

// CreateUserParams defines the parameters used to create a new user.
type CreateUserParams struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (x CreateUserParams) SpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("full_name", x.FullName),
		attribute.String("email", x.Email),
	}
}

// CreateUser creates a new user in the database.
func (x *SQL) CreateUser(ctx context.Context, params CreateUserParams) error {
	query := `
    INSERT INTO users (id, full_name, email, password_hash, created_at)
    VALUES ($1, $2, $3, $4, $5)
    `

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		log.Error().Err(err).Msg("db: error hashing password")
		return err
	}
	params.Password = string(passwordHash)

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
		if IsSentinelErr(err) {
			return SentinelErr(err)
		}
		log.Error().Err(err).Msg("db: error creating user")
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("db: error creating user, no rows affected")
		return err
	}
	if rowsAffected != 1 {
		log.Error().Err(err).Msg("db: error creating user, multiple rows affected")
		return ErrNoRowsAffected
	}

	return nil
}

// ReadUser reads a user from the database.
func (x *SQL) ReadUser(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User

	query := `
    SELECT id, full_name, email, password_hash, created_at, updated_at
    FROM users
    WHERE id = $1
    `

	if err := x.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		log.Error().Err(err).Msg("db: error reading user")
		return nil, err
	}

	return &user, nil
}

// ReadUserByEmail reads a user from the database by email.
func (x *SQL) ReadUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	query := `
    SELECT id, full_name, email, password_hash, created_at, updated_at
    FROM users
    WHERE email = $1
    `

	if err := x.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		log.Error().Err(err).Msg("db: error reading user")
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

// UpdateUser updates a user in the database.
func (x *SQL) UpdateUser(ctx context.Context, params UpdateUserParams) error {
	query := `
    UPDATE users
    SET full_name = $1, email = $2, updated_at = $3
    WHERE id = $4
    `

	res, err := x.db.ExecContext(
		ctx,
		query,
		params.FullName,
		params.Email,
		time.Now(),
		params.ID,
	)
	if err != nil {
		if IsSentinelErr(err) {
			return SentinelErr(err)
		}
		log.Error().Err(err).Msg("db: error updating user")
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("db: error updating user, no rows affected")
		return err
	}
	if rowsAffected != 1 {
		log.Error().Err(err).Msg("db: error updating user, multiple rows affected")
		return ErrNoRowsAffected
	}

	return nil
}

// DeleteUser deletes a user from the database.
func (x *SQL) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
    DELETE FROM users
    WHERE id = $1
    `

	result, err := x.db.ExecContext(ctx, query, id)
	if err != nil {
		if IsSentinelErr(err) {
			return SentinelErr(err)
		}
		log.Error().Err(err).Msg("db: error deleting user")
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		if rowsAffected > 1 {
			log.Error().Err(err).Msg("db: error deleting user, multiple rows affected")
			return err
		}
		return ErrUserNotFound
	}

	return nil
}
