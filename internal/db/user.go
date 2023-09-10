package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/Salam4nder/user/pkg/util"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// CreateUser creates a new user in the database as a transaction.
func (x *SQL) CreateUser(ctx context.Context, params CreateUserParams) (*User, error) {
	var user User

	query := `
    INSERT INTO users (full_name, email, password_hash, created_at)
    VALUES ($1, $2, $3, $4)
    RETURNING id, full_name, email, password_hash, created_at, updated_at
    `

	passwordHash, err := util.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	params.Password = passwordHash

	if err := x.db.QueryRowContext(
		ctx,
		query,
		params.FullName,
		params.Email,
		params.Password,
		params.CreatedAt,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if IsSentinelErr(err) {
			return nil, SentinelErr(err)
		}

		return nil, err
	}

	return &user, nil
}

// ReadUser reads a user from the database.
func (s *SQL) ReadUser(
	ctx context.Context,
	id uuid.UUID,
) (*User, error) {
	var user User

	query := `
    SELECT id, full_name, email, password_hash, created_at, updated_at
    FROM users
    WHERE id = $1
    `

	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

// ReadUserByEmail reads a user from the database by email.
func (s *SQL) ReadUserByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	var user User

	query := `
    SELECT id, full_name, email, password_hash, created_at, updated_at
    FROM users
    WHERE email = $1
    `

	if err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
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

// UpdateUserTx updates a user in the database as a transaction.
func (s *SQL) UpdateUserTx(
	ctx context.Context,
	params UpdateUserParams,
) (*User, error) {
	var user User

	query := `
    UPDATE users
    SET full_name = $1, email = $2, updated_at = $3
    WHERE id = $4
    RETURNING id, full_name, email, password_hash, created_at, updated_at
    `

	if err := s.execTx(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(
			ctx,
			query,
			params.FullName,
			params.Email,
			time.Now(),
			params.ID,
		).Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	}); err != nil {
		if pqErr, exists := err.(*pq.Error); exists {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, ErrDuplicateEmail

			default:
				return nil, err
			}
		} else if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

// DeleteUserTx deletes a user from the database.
func (s *SQL) DeleteUserTx(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
    DELETE FROM users
    WHERE id = $1
    `

	if err := s.execTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected < 1 {
			return ErrUserNotFound
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
