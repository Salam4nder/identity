package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/pkg/validation"
	"github.com/google/uuid"
	otelCode "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (x *UserServer) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*emptypb.Empty, error) {
	var err error
	ctx, span := tracer.Start(ctx, "rpc.CreateUser")
	defer func() {
		if err != nil {
			span.SetStatus(otelCode.Error, err.Error())
			span.RecordError(err)
		}
		span.End()
	}()
	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.Warn("CreateUser: GenSpanAttributes", "err", err.Error())
	}

	if err = validateCreateUserRequest(req); err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	params := db.CreateUserParams{
		ID:        uuid.New(),
		FullName:  req.GetFullName(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		CreatedAt: time.Now(),
	}
	if err = x.storage.CreateUser(ctx, params); err != nil {
		if errors.Is(err, db.ErrDuplicateEmail) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, internalServerError()
	}

	return &emptypb.Empty{}, nil
}

// validateCreateUserRequest returns nil if the request is valid.
func validateCreateUserRequest(req *gen.CreateUserRequest) error {
	if req == nil {
		return requestIsNilError()
	}
	var (
		fullNameErr error
		emailErr    error
		passwordErr error
	)

	if err := validation.FullName(req.GetFullName()); err != nil {
		fullNameErr = fmt.Errorf("grpc: full_name %w", err)
	}

	if err := validation.Email(req.GetEmail()); err != nil {
		emailErr = fmt.Errorf("grpc: email %w", err)
	}

	if err := validation.Password(req.GetPassword()); err != nil {
		passwordErr = fmt.Errorf("grpc: password %w", err)
	}

	return errors.Join(fullNameErr, emailErr, passwordErr)
}

func userToProtoResponse(user *db.User) *gen.UserResponse {
	return &gen.UserResponse{
		Id:        user.ID.String(),
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
