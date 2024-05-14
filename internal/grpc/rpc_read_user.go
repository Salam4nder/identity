package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/google/uuid"
	otelCode "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReadUser returns a user by ID.
func (x *UserServer) ReadUser(ctx context.Context, req *gen.UserID) (*gen.UserResponse, error) {
	var err error
	ctx, span := tracer.Start(ctx, "rpc.ReadUser")
	defer func() {
		if err != nil {
			span.SetStatus(otelCode.Error, err.Error())
			span.RecordError(err)
		}
		span.End()
	}()

	if req == nil {
		return nil, requestIsNilError()
	}

	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.Warn("ReadUser: GenSpanAttributes", "err", err)
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
		return nil, internalServerError()
	}

	return userToProtoResponse(user), nil
}
