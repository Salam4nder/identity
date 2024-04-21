package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/pkg/validation"
	"github.com/google/uuid"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelCode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	serverStr  string = "UserServer"
	handlerStr string = "CreateUser"
)

var tracer = otel.Tracer(serverStr)

func spanAttribures(req *gen.CreateUserRequest) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("full_name", req.FullName),
		attribute.String("email", req.Email),
		attribute.Int("password_length", len(req.Password)),
	}
}

func (x *UserServer) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*emptypb.Empty, error) {
	if req == nil {
		return &emptypb.Empty{}, requestIsNilError()
	}

	ctx, span := tracer.Start(ctx, handlerStr, trace.WithAttributes(spanAttribures(req)...))
	defer span.End()

	if err := validateCreateUserRequest(req); err != nil {
		span.SetStatus(otelCode.Error, err.Error())
		span.RecordError(err)
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	params := db.CreateUserParams{
		ID:        uuid.New(),
		FullName:  req.GetFullName(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		CreatedAt: time.Now(),
	}
	if err := x.storage.CreateUser(ctx, params); err != nil {
		span.SetStatus(otelCode.Error, err.Error())
		span.RecordError(err)
		if errors.Is(err, db.ErrDuplicateEmail) {
			return &emptypb.Empty{}, status.Error(codes.AlreadyExists, err.Error())
		}
		return &emptypb.Empty{}, internalServerError(err)
	}

	return &emptypb.Empty{}, nil
}

// validateCreateUserRequest returns nil if the request is valid.
func validateCreateUserRequest(req *gen.CreateUserRequest) error {
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
