package database

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type RowsAffectedError struct {
	want   int64
	actual int64
}

func (x RowsAffectedError) Error() string {
	return fmt.Sprintf("database: %d rows affected, expected %d", x.actual, x.want)
}

func NewRowsAffectedError(want, actual int64) RowsAffectedError {
	return RowsAffectedError{
		want:   want,
		actual: actual,
	}
}

// InputError is returned in case of input errors that are meant to be
// sent back to the user.
type InputError struct {
	field string
	value any
}

func (x InputError) Error() string {
	return fmt.Sprintf("database: input error on %s with value %v", x.field, x.value)
}

func NewInputError(field string, value any) InputError {
	return InputError{
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
}

func (x DuplicateEntryError) Error() string {
	return fmt.Sprintf("database: %s already exists", x.entity)
}

func NewDuplicateEntryError(entity string) DuplicateEntryError {
	return DuplicateEntryError{
		entity: entity,
	}
}

type NotFoundError struct {
	entity string
	id     any
}

func (x NotFoundError) Error() string {
	return fmt.Sprintf("database: %s not for id %v", x.entity, x.id)
}

func NewNotFoundError(entity string, id any) NotFoundError {
	return NotFoundError{
		entity: entity,
		id:     id,
	}
}
