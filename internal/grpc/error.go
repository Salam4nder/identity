package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func requestIsNilError() error {
	return status.Error(codes.InvalidArgument, "request is nil")
}

func invalidIDError() error {
	return status.Error(codes.InvalidArgument, "invalid id")
}

func internalServerError() error {
	return status.Error(codes.Internal, "internal server error")
}
