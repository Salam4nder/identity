package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ReadByEmail returns a user by email.
func (x *UserServer) ReadByEmail(ctx context.Context, req *gen.ReadByEmailRequest) (*gen.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "rpc.ReadByEmail")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError(span)
	}
	if req.GetEmail() == "" {
		return nil, invalidArgumentError(errors.New("email is required"), span, "email is required")
	}

	user, err := x.storage.ReadUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, notFoundError(err, span, "user not found")
		}
		return nil, internalServerError(err, span)
	}

	return &gen.UserResponse{
		Id:        user.ID.String(),
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}
