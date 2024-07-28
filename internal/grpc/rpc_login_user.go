package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	grpcUtil "github.com/Salam4nder/user/pkg/grpc"
	"github.com/Salam4nder/user/pkg/validation"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// LoginUser logs in a user and returns a session and a token payload.
func (x *UserServer) LoginUser(ctx context.Context, req *gen.LoginUserRequest) (*gen.LoginUserResponse, error) {
	ctx, span := tracer.Start(ctx, "rpc.LoginUser")
	defer span.End()

	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.Warn("LoginUser: GenSpanAttributes", "err", err)
	}

	if err = validateLoginUserRequest(req); err != nil {
		return nil, invalidArgumentError(err, span, err.Error())
	}

	user, err := x.storage.ReadUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, notFoundError(err, span, "user not found")
		}
		return nil, internalServerError(err, span)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, unauthenticatedError(err, span, "invalid password")
	}

	accessToken, accessPayload, err := x.tokenMaker.NewToken(
		req.GetEmail(),
		x.accessTokenDuration,
	)
	if err != nil {
		return nil, internalServerError(err, span)
	}

	refreshToken, refreshPayload, err := x.tokenMaker.NewToken(
		req.GetEmail(),
		x.refreshTokenDuration,
	)
	if err != nil {
		return nil, internalServerError(err, span)
	}

	metadata := grpcUtil.MetadataFromContext(ctx)

	if err = x.storage.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:        user.Email,
		ClientIP:     metadata.ClientIP,
		UserAgent:    metadata.UserAgent,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshPayload.ExpiresAt,
	}); err != nil {
		return nil, internalServerError(err, span)
	}

	// TODO(kg): Fix refresh token expiration.
	return &gen.LoginUserResponse{
		User: &gen.UserResponse{
			Id:        user.ID.String(),
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
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

	if err := validation.Email(req.GetEmail()); err != nil {
		return err
	}

	return nil
}
