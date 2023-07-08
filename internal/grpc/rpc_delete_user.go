package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/proto/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser deletes a user by id. Returns an error if the user couldn't be deleted,
// if the user doesn't exist or if the request is invalid.
func (s *userServer) DeleteUser(
	ctx context.Context, req *gen.UserID) (*emptypb.Empty, error) {
	if req == nil {
		return nil, requestIsNilError()
	}

	if req.GetId() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"ID can not be empty",
		)
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"ID is invalid",
		)
	}

	if err = s.storage.DeleteUserTx(ctx, id); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		default:
			s.logger.Error().Err(err).Msg("failed to delete user")

			return nil, internalServerError()
		}
	}

	return &emptypb.Empty{}, nil
}
