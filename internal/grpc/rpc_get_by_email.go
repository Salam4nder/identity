package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetByEmail returns a user by email. Returns an error if the user couldn't be not
// found or if the request is invalid.
func (s *userServer) GetByEmail(
	ctx context.Context, req *gen.GetByEmailRequest) (*gen.UserResponse, error) {
	if req == nil {
		return nil, requestIsNilError()
	}

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email can not be empty")
	}

	user, err := s.storage.ReadUserByEmail(ctx, req.GetEmail())
	if err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		default:
			s.logger.Error().Err(err).Msg("failed to get user by email")
			return nil, internalServerError()
		}
	}

	return userToProtoResponse(user), nil
}
