package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/proto/gen"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser deletes a user by ID.
func (s *UserServer) DeleteUser(
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

	if err = s.storage.DeleteUser(ctx, id); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		default:
			log.Error().Err(err).Msg("failed to delete user")

			return nil, internalServerError()
		}
	}

	return &emptypb.Empty{}, nil
}
