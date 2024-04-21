package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser deletes a user by ID.
func (x *UserServer) DeleteUser(ctx context.Context, req *gen.UserID) (*emptypb.Empty, error) {
	if req == nil {
		return nil, requestIsNilError()
	}

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ID can not be empty")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	if err = x.storage.DeleteUser(ctx, id); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		default:
			return nil, internalServerError(err)
		}
	}

	return &emptypb.Empty{}, nil
}
