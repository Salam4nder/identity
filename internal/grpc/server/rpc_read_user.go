package server

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/proto/gen"
	"github.com/google/uuid"
	otelCode "go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ReadUser returns a user by ID.
func (x *UserServer) ReadUser(ctx context.Context, req *gen.UserID) (*gen.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "rpc.ReadUser")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError(span)
	}

	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.Warn("ReadUser: GenSpanAttributes", "err", err)
	}

	if req.GetId() == "" {
		span.SetStatus(otelCode.Error, "ID is empty")
		return nil, invalidArgumentError(err, span, "ID is empty")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, invalidArgumentError(err, span, "ID is invalid")
	}

	user, err := db.ReadUser(ctx, x.db, id)
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
