package server

import (
	"context"
	"log/slog"

	"github.com/Salam4nder/user/internal/auth/strategy"
	"github.com/Salam4nder/user/internal/observability/metrics"
	"github.com/Salam4nder/user/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (x *Identity) Register(ctx context.Context, req *gen.Input) (*emptypb.Empty, error) {
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError(span)
	}

	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.WarnContext(ctx, "server: getting span attributes", "err", err)
	}

	switch t := x.strategy.(type) {
	case *strategy.Credentials:
		creds := req.GetCredentials()
		if err := t.IngestInput(ctx, creds.GetEmail(), creds.GetPassword()); err != nil {
			// TODO(kg): Errs.
			invalidArgumentError(err, span, "msg")
		}
		if err := t.Register(ctx); err != nil {
			// TODO(kg): Errs.
			invalidArgumentError(err, span, "msg")
		}
	case *strategy.NoOp:
		slog.InfoContext(ctx, "server: no op register")
	}

	metrics.UsersActive.Inc()
	metrics.UsersRegistered.Inc()

	return &emptypb.Empty{}, nil
}
