// response.go contains common response functions with tracing.
package grpc

import (
	otelCode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func internalServerError(err error, span trace.Span) error {
	span.SetStatus(otelCode.Error, err.Error())
	span.RecordError(err)
	return status.Error(codes.Internal, "internal server error occurred, provide the traceID to support")
}

func invalidArgumentError(err error, span trace.Span, msg string) error {
	span.SetStatus(otelCode.Error, err.Error())
	span.RecordError(err)
	return status.Error(codes.InvalidArgument, msg)
}

func alreadyExistsError(err error, span trace.Span, msg string) error {
	span.SetStatus(otelCode.Error, err.Error())
	span.RecordError(err)
	return status.Error(codes.AlreadyExists, msg)
}

func requestIsNilError(span trace.Span) error {
	span.SetStatus(otelCode.Error, "request is nil")
	return status.Error(codes.InvalidArgument, "request is nil")
}

func unauthenticatedError(err error, span trace.Span, msg string) error {
	span.SetStatus(otelCode.Error, err.Error())
	span.RecordError(err)
	return status.Error(codes.Unauthenticated, msg)
}

// nolint
func notFoundError(err error, span trace.Span, msg string) error {
	span.SetStatus(otelCode.Error, err.Error())
	span.RecordError(err)
	return status.Error(codes.NotFound, msg)
}
