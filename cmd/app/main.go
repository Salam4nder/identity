package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/db"
	internalGRPC "github.com/Salam4nder/user/internal/grpc"
	"github.com/Salam4nder/user/internal/proto/gen"
	grpcUtil "github.com/Salam4nder/user/pkg/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stimtech/go-migration"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	// PingTimeout is the maximum duration for waiting on ping.
	PingTimeout = 5 * time.Second
	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	ReadTimeout = 10 * time.Second
	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read.
	WriteTimeout = 10 * time.Second
	// EnvironmentDev is the development environment.
	EnvironmentDev = "dev"
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		PingTimeout,
	)
	defer cancel()

	cfg, err := config.New()
	exitWithError(err)

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.Environment == EnvironmentDev {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	sql, err := db.NewSQLDatabase(ctx, cfg.PSQL)
	exitWithError(err)

	migration := migration.New(sql.DB(), zap.NewNop())
	if err = migration.Migrate(); err != nil {
		exitWithError(err)
	}

	userServer, err := internalGRPC.NewUserServer(sql, cfg.UserService)
	exitWithError(err)

	grpcListener, err := net.Listen("tcp", cfg.Server.GRPCAddr())
	exitWithError(err)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcUtil.LoggerInterceptor,
			recovery.UnaryServerInterceptor(),
		),
	)
	gen.RegisterUserServer(grpcServer, userServer)
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			exitWithError(err)
		}
	}()

	log.Info().
		Str("address", cfg.Server.GRPCAddr()).
		Msg("gRPC server is running")

	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err = gen.RegisterUserHandlerFromEndpoint(
		ctx,
		mux,
		cfg.Server.GRPCAddr(),
		dialOpts,
	); err != nil {
		exitWithError(err)
	}

	server := &http.Server{
		Handler:      mux,
		Addr:         cfg.Server.HTTPAddr(),
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			exitWithError(err)
		}
	}()

	log.Info().
		Str("address", cfg.Server.HTTPAddr()).
		Msg("gRPC gateway is running")

	// Wait for interrupt signal to gracefully shutdown the server with.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")
	grpcServer.GracefulStop()
	if err := server.Shutdown(ctx); err != nil {
		exitWithError(err)
	}
	log.Info().Msg("server gracefully stopped")
}

func exitWithError(err error) {
	if err != nil {
		log.Error().Err(err).Msg("main: exit with error")
		os.Exit(1)
	}
}
