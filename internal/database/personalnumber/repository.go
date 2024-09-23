package personalnumber

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Salam4nder/identity/internal/database"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("personal_number")

const Tablename = "personal_numbers"

func Insert(ctx context.Context, db *sql.DB, id uint64) error {
	ctx, span := tracer.Start(ctx, "Insert")
	defer span.End()
	span.SetAttributes(attribute.Int64("id", int64(id)))

	query := `INSERT INTO personal_numbers (id) VALUES ($1)`
	span.SetAttributes(attribute.String("query", query))

	res, err := db.ExecContext(ctx, query, id)
	if err != nil {
		if database.IsPSQLDuplicateEntryError(err) {
			return database.NewDuplicateEntryError(ctx, err, "personal_number")
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

func Get(ctx context.Context, db *sql.DB, id uint64) (uint64, error) {
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()
	span.SetAttributes(attribute.Int64("id", int64(id)))

	if id == 0 {
		return 0, database.NewInputError(ctx, errors.New("id is empty"), "id", id)
	}

	query := `
    SELECT id personal_numbers
    WHERE id = $1
    `
	span.SetAttributes(attribute.String("query", query))

	var res uint64
	if err := db.QueryRowContext(ctx, query, id).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, database.NewNotFoundError(ctx, err, "personal_number", id)
		}
		return 0, database.NewOperationFailedError(ctx, err)
	}
	return res, nil
}

func Delete(ctx context.Context, db *sql.DB, id uint64) error {
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()
	span.SetAttributes(attribute.Int64("id", int64(id)))

	query := `
    DELETE FROM personal_numbers
    WHERE id = $1
    `
	span.SetAttributes(attribute.String("query", query))

	res, err := db.ExecContext(
		ctx,
		query,
		id,
	)
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
