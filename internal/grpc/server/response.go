package server

import (
	"context"

	otelCode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func requestIsNilError() error {
	return status.Error(codes.InvalidArgument, "request is nil")
}

func internalServerError(ctx context.Context, err error) error {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(otelCode.Error, err.Error())
		span.RecordError(err)
	}
	return status.Error(codes.Internal, "internal server error, please provide the traceID to support")
}

func invalidArgumentError(ctx context.Context, err error, msg string) error {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(otelCode.Error, err.Error())
		span.RecordError(err)
	}
	return status.Error(codes.InvalidArgument, msg)
}

func alreadyExistsError(ctx context.Context, err error, msg string) error {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(otelCode.Error, err.Error())
		span.RecordError(err)
	}
	return status.Error(codes.AlreadyExists, msg)
}

// func unauthenticatedError(err error, span trace.Span, msg string) error {
// 	span.SetStatus(otelCode.Error, err.Error())
// 	span.RecordError(err)
// 	return status.Error(codes.Unauthenticated, msg)
// }

// func notFoundError(ctx context.Context, err error, msg string) error {
// 	if err != nil {
// 		span := trace.SpanFromContext(ctx)
// 		span.SetStatus(otelCode.Error, err.Error())
// 		span.RecordError(err)
// 	}
// 	return status.Error(codes.NotFound, msg)
// }
