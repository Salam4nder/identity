package personalnumber

import (
	"context"
	"database/sql"
	"time"

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

	query := `INSERT INTO personal_numbers (id, created_at)
    VALUES ($1, $2)
    `
	span.SetAttributes(attribute.String("query", query))

	res, err := db.ExecContext(
		ctx,
		query,
		id,
		time.Now(),
	)
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
