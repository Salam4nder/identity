package grpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/pkg/validation"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (x *UserServer) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*emptypb.Empty, error) {
	ctx, span := tracer.Start(ctx, "rpc.CreateUser")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError(span)
	}

	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.Warn("grpc.CreateUser: GenSpanAttributes", "err", err.Error())
	}

	if err = validateCreateUserRequest(req); err != nil {
		return nil, invalidArgumentError(err, span, err.Error())
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
			return nil, alreadyExistsError(err, span, "user with the provided email already exists")
		}
		return nil, internalServerError(err, span)
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
		fullNameErr = err
	}

	if err := validation.Email(req.GetEmail()); err != nil {
		emailErr = err
	}

	if err := validation.Password(req.GetPassword()); err != nil {
		passwordErr = err
	}

	return errors.Join(fullNameErr, emailErr, passwordErr)
}
