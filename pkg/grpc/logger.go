package grpc

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggerInterceptor logs gRPC requests.
func LoggerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()

	result, err := handler(ctx, req)

	code := codes.Unknown

	if status, exists := status.FromError(err); exists {
		code = status.Code()
	}

	duration := time.Since(startTime)

	log.Info().
		Str("method", info.FullMethod).
		Int("status_code", int(code)).
		Str("status_text", code.String()).
		Dur("duration", duration).
		Msg("gRPC request received")

	return result, err
}
