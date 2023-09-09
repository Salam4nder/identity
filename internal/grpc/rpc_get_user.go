package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/proto/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUser returns a user by id. Returns an error if the user couldn't be found
// or if the request is invalid.
func (s *UserServer) GetUser(
	ctx context.Context, req *gen.UserID) (*gen.UserResponse, error) {
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

	user, err := s.storage.ReadUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		default:
			s.logger.Error().Err(err).Msg("failed to get user")
			return nil, internalServerError()
		}
	}

	return userToProtoResponse(user), nil
}
