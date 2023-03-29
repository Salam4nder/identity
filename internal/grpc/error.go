package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func requestIsNilError() error {
	return status.Error(codes.InvalidArgument, "request is nil")
}

func internalServerError() error {
	return status.Error(codes.Internal, "internal server error")
}

func unauthenticatedError() error {
	return status.Error(codes.Unauthenticated, "unauthenticated")
}
