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
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser deletes a user by ID.
func (x *UserServer) DeleteUser(ctx context.Context, req *gen.UserID) (*emptypb.Empty, error) {
	var err error
	ctx, span := tracer.Start(ctx, "rpc.DeleteUser")
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
		slog.Warn("DeleteUser: GenSpanAttributes", "err", err)
	}

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ID can not be empty")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	if err = x.storage.DeleteUser(ctx, id); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		default:
			return nil, internalServerError()
		}
	}

	return &emptypb.Empty{}, nil
}
