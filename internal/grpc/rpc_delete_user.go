package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/internal/observability/metrics"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser deletes a user by ID.
func (x *UserServer) DeleteUser(ctx context.Context, req *gen.UserID) (*emptypb.Empty, error) {
	ctx, span := tracer.Start(ctx, "rpc.DeleteUser")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError(span)
	}

	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.Warn("DeleteUser: GenSpanAttributes", "err", err)
	}

	if req.GetId() == "" {
		return nil, invalidArgumentError(errors.New("ID is required"), span, "ID is required")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, invalidArgumentError(err, span, "invalid ID")
	}

	if err = x.storage.DeleteUser(ctx, id); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			return nil, notFoundError(err, span, "user not found")

		default:
			return nil, internalServerError(err, span)
		}
	}

	metrics.UsersActive.Dec()

	return &emptypb.Empty{}, nil
}
