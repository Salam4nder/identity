package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	grpcUtil "github.com/Salam4nder/user/pkg/grpc"
	"github.com/Salam4nder/user/pkg/validation"
	otelCode "go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// LoginUser logs in a user and returns a session and a token payload.
func (x *UserServer) LoginUser(ctx context.Context, req *gen.LoginUserRequest) (*gen.LoginUserResponse, error) {
	var err error
	ctx, span := tracer.Start(ctx, "rpc.LoginUser")
	defer func() {
		if err != nil {
			span.SetStatus(otelCode.Error, err.Error())
			span.RecordError(err)
		}
		span.End()
	}()

	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.Warn("LoginUser: GenSpanAttributes", "err", err)
	}

	if err = validateLoginUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := x.storage.ReadUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found.")
		}
		return nil, internalServerError()
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	accessToken, accessPayload, err := x.tokenMaker.NewToken(
		req.GetEmail(),
		x.accessTokenDuration,
	)
	if err != nil {
		return nil, internalServerError()
	}

	refreshToken, refreshPayload, err := x.tokenMaker.NewToken(
		req.GetEmail(),
		x.refreshTokenDuration,
	)
	if err != nil {
		return nil, internalServerError()
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
		return nil, internalServerError()
	}

	// reminder to fix expiration timing on refresh token
	return &gen.LoginUserResponse{
		User:                  userToProtoResponse(user),
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

	if err := validation.Email(req.GetEmail()); err != nil {
		emailErr = err
	}

	if err := validation.Password(req.GetPassword()); err != nil {
		passwordErr = err
	}

	return errors.Join(emailErr, passwordErr)
}
