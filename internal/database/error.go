package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var ErrUnexpectedRowsAffectedError = errors.New("database: unexpected rows affected")

type RowsAffectedError struct {
	want   int64
	actual int64
}

func (x RowsAffectedError) Error() string {
	return fmt.Sprintf("database: %d rows affected, expected %d", x.actual, x.want)
}

func NewRowsAffectedError(ctx context.Context, err error, want, actual int64) RowsAffectedError {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	return RowsAffectedError{
		want:   want,
		actual: actual,
	}
}

// InputError is returned in case of input errors that are meant to be
// sent back to the user.
type InputError struct {
	inner error
	field string
	value any
}

func (x InputError) Error() string {
	return fmt.Sprintf("database: input error on %s with value %v", x.field, x.value)
}

func (x InputError) Inner() string {
	if x.inner == nil {
		return ""
	}
	return x.inner.Error()
}

func NewInputError(ctx context.Context, err error, field string, value any) InputError {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	return InputError{
		inner: err,
		field: field,
		value: value,
	}
}

func IsPSQLDuplicateEntryError(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation"
}

type DuplicateEntryError struct {
	entity string
	inner  error
}

func (x DuplicateEntryError) Error() string {
	return fmt.Sprintf("database: %s already exists", x.entity)
}

func (x DuplicateEntryError) Inner() string {
	if x.inner == nil {
		return ""
	}
	return x.inner.Error()
}

func NewDuplicateEntryError(ctx context.Context, err error, entity string) DuplicateEntryError {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	return DuplicateEntryError{
		inner:  err,
		entity: entity,
	}
}

type NotFoundError struct {
	id     any
	inner  error
	entity string
}

func (x NotFoundError) Error() string {
	return fmt.Sprintf("database: %s not for id %v", x.entity, x.id)
}

func (x NotFoundError) Inner() string {
	if x.inner == nil {
		return ""
	}
	return x.inner.Error()
}

func NewNotFoundError(ctx context.Context, err error, entity string, id any) NotFoundError {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	return NotFoundError{
		id:     id,
		inner:  err,
		entity: entity,
	}
}

type OperationFailedError struct {
	inner error
}

func (x OperationFailedError) Error() string {
	if x.inner == nil {
		return ""
	}
	return x.inner.Error()
}

func NewOperationFailedError(ctx context.Context, err error) OperationFailedError {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	return OperationFailedError{
		inner: err,
	}
}
