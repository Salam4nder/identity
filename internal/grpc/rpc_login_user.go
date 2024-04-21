package grpc

import (
	"context"
	"errors"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	grpcUtil "github.com/Salam4nder/user/pkg/grpc"
	"github.com/Salam4nder/user/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// LoginUser logs in a user and returns a session and a token payload.
func (x *UserServer) LoginUser(
	ctx context.Context,
	req *gen.LoginUserRequest,
) (*gen.LoginUserResponse, error) {
	if err := validateLoginUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := x.storage.ReadUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "handlers: user not found")
		}
		return nil, internalServerError()
	}

	if err := util.ComparePasswordHash(req.Password, user.PasswordHash); err != nil {
		return nil, status.Error(codes.Unauthenticated, "handlers: invalid password")
	}

	accessToken, accessPayload, err := x.tokenMaker.NewToken(
		req.GetEmail(),
		x.accessTokenDuration,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refreshToken, refreshPayload, err := x.tokenMaker.NewToken(
		req.GetEmail(),
		x.refreshTokenDuration,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	metadata := grpcUtil.MetadataFromContext(ctx)

	session, err := x.storage.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:        user.Email,
		ClientIP:     metadata.ClientIP,
		UserAgent:    metadata.UserAgent,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})
	if err != nil {
		return nil, internalServerError()
	}

	// reminder to fix expiration timing on refresh token
	return &gen.LoginUserResponse{
		User:                  userToProtoResponse(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt),
	}, nil
}

func validateLoginUserRequest(req *gen.LoginUserRequest) error {
	if req == nil {
		return errors.New("request can not be nil")
	}

	var (
		emailErr    error
		passwordErr error
	)

	if err := util.ValidateEmail(req.GetEmail()); err != nil {
		emailErr = err
	}

	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		passwordErr = err
	}

	return errors.Join(emailErr, passwordErr)
}
