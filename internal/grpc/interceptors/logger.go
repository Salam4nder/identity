package interceptors

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tracer = otel.Tracer("LoggingInterceptor")

// UnaryLoggerInterceptor logs gRPC requests.
func UnaryLoggerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	ctx, span := tracer.Start(ctx, info.FullMethod)
	defer span.End()

	startTime := time.Now()

	result, err := handler(ctx, req)

	duration := time.Since(startTime)
	code := codes.Unknown
	if status, exists := status.FromError(err); exists {
		code = status.Code()
	}

	span.SetAttributes(
		attribute.String("protocol", "grpc"),
		attribute.String("method", info.FullMethod),
		attribute.Int("status_code", int(code)),
		attribute.String("status_text", code.String()),
		attribute.Int64("duration", duration.Microseconds()),
	)

	attrs := []slog.Attr{
		slog.String("protocol", "grpc"),
		slog.String("method", info.FullMethod),
		slog.Int("status_code", int(code)),
		slog.String("status_text", code.String()),
		slog.Duration("duration", time.Duration(duration.Milliseconds())),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "log interceptor:", attrs...)
	} else {
		slog.LogAttrs(ctx, slog.LevelInfo, "log interceptor:", attrs...)
	}

	return result, err
}
