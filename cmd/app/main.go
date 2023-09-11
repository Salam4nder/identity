package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
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
	cfg, err := config.New()
	fatalExitOnErr(err)

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.Environment == EnvironmentDev {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		PingTimeout,
	)
	defer cancel()

	sql, err := db.NewSQLDatabase(ctx, cfg.PSQL)
	fatalExitOnErr(err)

	migration := migration.New(sql.DB(), zap.NewNop())
	if err := migration.Migrate(); err != nil {
		fatalExitOnErr(err)
	}

	userServer, err := internalGRPC.NewUserServer(sql, cfg.Service)
	fatalExitOnErr(err)

	grpcListener, err := net.Listen("tcp", cfg.Server.GRPCAddr())
	fatalExitOnErr(err)

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
			fatalExitOnErr(err)
		}
	}()

	log.Info().
		Str("address", cfg.Server.GRPCAddr()).
		Msg("gRPC server is running")

	mux := runtime.NewServeMux()

	// dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := gen.RegisterUserHandlerFromEndpoint(
		ctx,
		mux,
		cfg.Server.GRPCAddr(),
		dialOpts,
	); err != nil {
		fatalExitOnErr(err)
	}

	server := &http.Server{
		Handler:      mux,
		Addr:         cfg.Server.HTTPAddr(),
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fatalExitOnErr(err)
		}
	}()

	log.Info().
		Str("address", cfg.Server.HTTPAddr()).
		Msg("gRPC gateway is running")

	// Wait for interrupt or kill signal to gracefully shutdown the server with.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	log.Info().Msg("shutting down server...")
	grpcServer.GracefulStop()
	server.Shutdown(ctx)
	log.Info().Msg("server gracefully stopped")
}

func fatalExitOnErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("main fatal exit: failed to start user service")
	}
}
