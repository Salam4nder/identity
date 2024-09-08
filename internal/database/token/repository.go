package token

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Salam4nder/identity/internal/database"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("token")

const Tablename = "tokens"

func Insert(ctx context.Context, db *sql.DB, token string) error {
	ctx, span := tracer.Start(ctx, "Insert")
	defer span.End()
	span.SetAttributes(attribute.String("token", token))

	if token == "" {
		return database.NewInputError(
			ctx,
			errors.New("token: token is empty"),
			"token",
			token,
		)
	}

	query := `INSERT INTO tokens (token) VALUES ($1)`
	span.SetAttributes(attribute.String("query", query))
	res, err := db.ExecContext(ctx, query, token)
	if err != nil {
		if database.IsPSQLDuplicateEntryError(err) {
			return database.NewDuplicateEntryError(ctx, err, "tokens")
		}
		return database.NewOperationFailedError(ctx, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return database.NewOperationFailedError(ctx, err)
	}
	if rowsAffected != 1 {
		return database.NewRowsAffectedError(ctx, database.ErrUnexpectedRowsAffectedError, 1, rowsAffected)
	}

	return nil
}

func Get(ctx context.Context, db *sql.DB, token string) (string, error) {
	ctx, span := tracer.Start(ctx, "Get")
	defer span.End()
	span.SetAttributes(attribute.String("token", token))

	if token == "" {
		return "", database.NewInputError(
			ctx,
			errors.New("token: token is empty"),
			"token",
			token,
		)
	}

	var s string
	query := `SELECT token FROM tokens WHERE token = $1`
	span.SetAttributes(attribute.String("query", query))
	if err := db.QueryRowContext(ctx, query, token).Scan(&s); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", database.NewNotFoundError(ctx, err, "token", token)
		}
		return "", database.NewOperationFailedError(ctx, err)
	}

	return s, nil
}

func Delete(ctx context.Context, db *sql.DB, token string) error {
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()
	span.SetAttributes(attribute.String("token", token))

	if token == "" {
		return database.NewInputError(
			ctx,
			errors.New("token: token is empty"),
			"token",
			token,
		)
	}

	query := `DELETE FROM tokens WHERE token = $1`
	span.SetAttributes(attribute.String("query", query))
	res, err := db.ExecContext(ctx, query, token)
	if err != nil {
		return database.NewOperationFailedError(ctx, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return database.NewOperationFailedError(ctx, err)
	}
	if rowsAffected != 1 {
		return database.NewRowsAffectedError(ctx, database.ErrUnexpectedRowsAffectedError, 1, rowsAffected)
	}

	return nil
}
