package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/pkg/validation"
	"github.com/google/uuid"
	otelCode "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser updates a user by ID.
func (x *UserServer) UpdateUser(ctx context.Context, req *gen.UpdateUserRequest) (*emptypb.Empty, error) {
	var err error
	ctx, span := tracer.Start(ctx, "rpc.UpdateUser")
	defer func() {
		if err != nil {
			span.SetStatus(otelCode.Error, err.Error())
			span.RecordError(err)
		}
		span.End()
	}()

	if req == nil {
		return &emptypb.Empty{}, requestIsNilError()
	}

	atts, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(atts...)
	} else {
		slog.Warn("UpdateUser: GenSpanAttributes", "err", err)
	}

	authPayload, err := x.authorizeUser(ctx)
	if err != nil {
		return &emptypb.Empty{}, unauthenticatedError(err)
	}

	if err = validateUpdateUserRequest(req); err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if authPayload.Email != req.GetEmail() {
		return &emptypb.Empty{}, status.Errorf(
			codes.PermissionDenied,
			"not owner of provided email",
		)
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	params := db.UpdateUserParams{
		ID:       id,
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	}

	if err = x.storage.UpdateUser(ctx, params); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())

		default:
			return &emptypb.Empty{}, internalServerError()
		}
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
