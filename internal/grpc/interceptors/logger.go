package interceptors

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tracer = otel.Tracer("LoggingInterceptor")

// UnaryLoggerInterceptor logs gRPC requests.
func UnaryLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	ctx, span := tracer.Start(ctx, info.FullMethod)
	defer span.End()
	startTime := time.Now()

	result, err := handler(ctx, req)
	traceID := span.SpanContext().TraceID().String()

	code := codes.Unknown
	if status, exists := status.FromError(err); exists {
		code = status.Code()
	}

	logger := log.Info()
	if err != nil {
		span.RecordError(err)
		logger = log.Error().Err(err)
	}

	duration := time.Since(startTime)

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(code)).
		Str("status_text", code.String()).
		Dur("duration", duration).
		Str("trace_id", traceID).
		Send()

	return result, err
}
