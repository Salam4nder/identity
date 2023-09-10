package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/proto/gen"
	grpcutil "github.com/Salam4nder/user/pkg/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is the interface for the gRPC server.
type Server interface {
	ServeGRPC() error
	ServeGRPCGateway() error
}

type server struct {
	userSrvc *UserServer
	cfg      *config.Server
}

// NewServer creates new gRPC server.
func NewServer(srvc *UserServer, cfg *config.Server) Server {
	return &server{
		userSrvc: srvc,
		cfg:      cfg,
	}
}

// ServeGRPC starts the gRPC server.
func (s *server) ServeGRPC() error {
	listener, err := net.Listen("tcp", s.cfg.GRPCAddr())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcutil.LoggerInterceptor,
			recovery.UnaryServerInterceptor(),
		),
	)

	gen.RegisterUserServer(grpcServer, s.userSrvc)
	reflection.Register(grpcServer)

	log.Info().
		Str("address", s.cfg.GRPCAddr()).
		Msg("gRPC server is running")

	return grpcServer.Serve(listener)
}

// ServeGRPCGateway starts the gRPC gateway.
func (s *server) ServeGRPCGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	err := gen.RegisterUserHandlerFromEndpoint(ctx, mux, s.cfg.GRPCAddr(), dialOpts)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	server := &http.Server{
		Handler: mux,
	}

	listener, err := net.Listen("tcp", s.cfg.HTTPAddr())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Info().
		Str("address", s.cfg.HTTPAddr()).
		Msg("gRPC gateway is running")

	return server.Serve(listener)
}
