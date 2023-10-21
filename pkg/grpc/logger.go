package grpc

import (
	"context"
	"net/http"
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

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	duration := time.Since(startTime)

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(code)).
		Str("status_text", code.String()).
		Dur("duration", duration).
		Send()

	return result, err
}

// ResponseRecorder is a wrapper around http.ResponseWriter that records its
// HTTP status code and body size.
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

// WriteHeader records the HTTP status code.
func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Write is a wrapper around http.ResponseWriter.Write that records the
// response body.
func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

// HTTPLogger logs HTTP requests. Used for gRPC Gateway.
func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()

		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, req)

		duration := time.Since(startTime)

		logger := log.Info()

		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http 1.1").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Send()
	})
}
