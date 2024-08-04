package server

import (
	"context"
	"fmt"
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

	switch t := x.strategy.(type) {
	case *strategy.Credentials:
		attrs, err := GenSpanAttributes(req.GetCredentials())
		if err == nil {
			span.SetAttributes(attrs...)
		} else {
			slog.WarnContext(ctx, "server: getting span attributes", "err", err)
		}
		if err = t.IngestInput(ctx, strategy.CredentialsInput{
			Email:    req.GetCredentials().GetEmail(),
			Password: req.GetCredentials().GetPassword(),
		}); err != nil {
			return nil, invalidArgumentError(err, span, err.Error())
		}
		if err := t.Register(ctx); err != nil {
			// TODO(kg): Errs.
			return nil, invalidArgumentError(err, span, err.Error())
		}
	case *strategy.NoOp:
		slog.InfoContext(ctx, "server: no op register")
	default:
		slog.ErrorContext(ctx, fmt.Sprintf("server: unsupported strategy %T,", t))
		internalServerError(nil, span)
	}

	metrics.UsersActive.Inc()
	metrics.UsersRegistered.Inc()

	return &emptypb.Empty{}, nil
}
