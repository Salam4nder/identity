package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReadUser returns a user by ID.
func (x *UserServer) ReadUser(ctx context.Context, req *gen.UserID) (*gen.UserResponse, error) {
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

	user, err := x.storage.ReadUser(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Error().Err(err).Msg("grpc: failed to read user")

		return nil, internalServerError()
	}

	return userToProtoResponse(user), nil
}
