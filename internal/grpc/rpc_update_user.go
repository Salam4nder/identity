package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/pkg/validation"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser updates a user by ID.
func (x *UserServer) UpdateUser(ctx context.Context, req *gen.UpdateUserRequest) error {
	if req == nil {
		return requestIsNilError()
	}

	authPayload, err := x.authorizeUser(ctx)
	if err != nil {
		return unauthenticatedError(err)
	}

	if err = validateUpdateUserRequest(req); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if authPayload.Email != req.GetEmail() {
		return status.Errorf(
			codes.PermissionDenied,
			"not owner of provided email",
		)
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return status.Error(codes.InvalidArgument, "ID is invalid")
	}

	params := db.UpdateUserParams{
		ID:       id,
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	}

	if err = x.storage.UpdateUser(ctx, params); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return status.Error(codes.NotFound, err.Error())

		default:
			return internalServerError(err)
		}
	}
	return nil
}

// validateUpdateUserRequest returns nil if the request is valid.
func validateUpdateUserRequest(req *gen.UpdateUserRequest) error {
	if req.Id == "" {
		return errors.New("ID can not be empty")
	}

	var (
		fullNameErr error
		emailErr    error
	)

	if req.GetFullName() != "" {
		if err := validation.FullName(req.GetFullName()); err != nil {
			fullNameErr = err
		}
	}

	if req.GetEmail() != "" {
		if err := validation.Email(req.GetEmail()); err != nil {
			emailErr = err
		}
	}

	return errors.Join(fullNameErr, emailErr)
}
