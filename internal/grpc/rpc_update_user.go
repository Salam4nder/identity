package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/pkg/validation"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser updates a user by ID.
func (x *UserServer) UpdateUser(ctx context.Context, req *gen.UpdateUserRequest) (*emptypb.Empty, error) {
	ctx, span := tracer.Start(ctx, "rpc.UpdateUser")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError(span)
	}

	atts, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(atts...)
	} else {
		slog.Warn("UpdateUser: GenSpanAttributes", "err", err)
	}

	authPayload, err := x.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err, span, "no permission to access rpc")
	}

	if err = validateUpdateUserRequest(req); err != nil {
		return nil, invalidArgumentError(err, span, err.Error())
	}

	if authPayload.Email != req.GetEmail() {
		return nil, unauthenticatedError(errors.New("email does not match"), span, "no permission to access rpc")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, invalidArgumentError(err, span, "invalid ID")
	}

	params := db.UpdateUserParams{
		ID:       id,
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	}
	if err = x.storage.UpdateUser(ctx, params); err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, notFoundError(err, span, "user not found")
		}
		return nil, internalServerError(err, span)
	}
	return &emptypb.Empty{}, nil
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
