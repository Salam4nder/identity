package grpc

import (
	"fmt"
	"net"

	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/proto/pb"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	srvc   *userService
	cfg    *config.GRPCServer
	logger *zap.Logger
}

// NewServer creates new gRPC server.
func NewServer(
	srvc *userService,
	cfg *config.GRPCServer,
	logger *zap.Logger,
) *server {
	return &server{
		srvc:   srvc,
		cfg:    cfg,
		logger: logger,
	}
}

// Serve starts the gRPC server.
func (s *server) Serve() error {
	listener, err := net.Listen("tcp", s.cfg.Addr())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, s.srvc)
	reflection.Register(grpcServer)

	s.logger.Info(
		"gRPC server is running", zap.String("address", s.cfg.Addr()))

	return grpcServer.Serve(listener)
}
