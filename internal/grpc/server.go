package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/proto/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	userSrvc *userServer
	cfg      *config.Server
	logger   *zap.Logger
}

// NewServer creates new gRPC server.
func NewServer(
	srvc *userServer,
	cfg *config.Server,
	logger *zap.Logger,
) *server {
	return &server{
		userSrvc: srvc,
		cfg:      cfg,
		logger:   logger,
	}
}

// ServeGRPC starts the gRPC server.
func (s *server) ServeGRPC() error {
	listener, err := net.Listen("tcp", s.cfg.GRPCAddr())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	gen.RegisterUserServer(grpcServer, s.userSrvc)
	reflection.Register(grpcServer)

	s.logger.Info(
		"gRPC server is running", zap.String("address", s.cfg.GRPCAddr()))

	return grpcServer.Serve(listener)
}

// ServeGRPCGateway starts the gRPC gateway.
func (s *server) ServeGRPCGateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gen.RegisterUserHandlerFromEndpoint(ctx, mux, s.cfg.GRPCAddr(), opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	server := &http.Server{
		Handler: mux,
	}

	listener, err := net.Listen("tcp", s.cfg.HTTPAddr())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.logger.Info(
		"gRPC gateway is running", zap.String("address", s.cfg.HTTPAddr()))

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
