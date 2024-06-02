package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Salam4nder/user/internal/config"
	internalDB "github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/email"
	"github.com/Salam4nder/user/internal/event"
	internalGRPC "github.com/Salam4nder/user/internal/grpc"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/internal/grpc/interceptors"
	"github.com/Salam4nder/user/internal/otel"
	"github.com/Salam4nder/user/pkg/logger"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/nats-io/nats.go"
	"github.com/stimtech/go-migration"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgen "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	// TODO(kg): Move most of these to config.

	// serviceName is the name of the service.
	serviceName string = "user"
	// serviceVersion is the version of the service.
	serviceVersion string = "1.0.0"
	// migrationFolder is the folder where the migration files are stored.
	migrationFolder = "internal/db/migrations"
	// accessTokenDuration is the duration for which the access token is valid.
	accessTokenDuration = 15 * time.Minute
	// refreshTokenDuration is the duration for which the refresh token is valid.
	refreshTokenDuration = 7 * 24 * time.Hour
	// natsTimeout is the timeout for the NATS connection.
	natsTimeout = 5 * time.Second
)

// serviceID is the unique identifier of the service.
var serviceID string = uuid.New().String()

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	exitOnError(ctx, err)

	if cfg.Environment == "dev" {
		slog.SetDefault(slog.New(logger.NewOtelHandler(logger.NewTintHandler(os.Stdout, nil))))
	} else {
		slog.SetDefault(slog.New(logger.NewOtelHandler(slog.NewJSONHandler(os.Stdout, nil))))
	}

	otelShutdown, err := otel.Setup(ctx, otel.SetupOpts{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		ServiceID:      serviceID,
	})
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	db, err := sql.Open(cfg.PSQL.Driver(), cfg.PSQL.Addr())
	exitOnError(ctx, err)
	storage := internalDB.New(db)
	if err = storage.PingContext(ctx, 5 /*max tries*/); err != nil {
		exitOnError(ctx, err)
	}
	migration := migration.New(db, zap.NewNop()).WithFolder(migrationFolder)
	if err = migration.Migrate(); err != nil {
		exitOnError(ctx, err)
	}

	// NATS.
	natsClient, err := nats.Connect(
		cfg.NATS.Addr(),
		nats.Timeout(natsTimeout),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(20),
	)
	exitOnError(ctx, err)
	natsChan := make(chan *nats.Msg, 64)
	userSub, err := natsClient.ChanSubscribe(event.UserRegistered, natsChan)
	exitOnError(ctx, err)

	worker := event.NewWorker(email.NewNoOpSender())
	go worker.Work(ctx, natsChan)

	grpcListener, err := net.Listen("tcp", cfg.Server.GRPCAddr())
	exitOnError(ctx, err)
	grpcServer := grpc.NewServer(
		// grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			interceptors.UnaryLoggerInterceptor,
		),
	)
	healthServer := health.NewServer()
	healthgen.RegisterHealthServer(grpcServer, healthServer)
	userServer, err := internalGRPC.NewUserServer(
		storage,
		healthServer,
		natsClient,
		cfg.Server.SymmetricKey,
		accessTokenDuration,
		refreshTokenDuration,
	)
	exitOnError(ctx, err)
	gen.RegisterUserServer(grpcServer, userServer)
	reflection.Register(grpcServer)

	go userServer.MonitorHealth(ctx)

	srvErrChan := make(chan error, 1)
	go func() {
		srvErrChan <- grpcServer.Serve(grpcListener)
	}()
	slog.InfoContext(ctx, "main: gRPC server is running", "address", cfg.Server.GRPCAddr())

	select {
	case err := <-srvErrChan:
		slog.ErrorContext(ctx, "gRPC server error", "error", err)
	case <-ctx.Done():
		slog.InfoContext(ctx, "main: context done, shutting down...")
	}
	grpcServer.GracefulStop()
	err = errors.Join(err, db.Close())
	err = errors.Join(err, userSub.Unsubscribe())
	natsClient.Close()

	if err != nil {
		slog.ErrorContext(ctx, "main: error upon exit", "error", err)
	}
}

func exitOnError(ctx context.Context, err error) {
	if err != nil {
		slog.ErrorContext(ctx, "main: exit on error", "error", err)
		ctx.Done()
		os.Exit(1)
	}
}
