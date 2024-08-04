package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Salam4nder/user/internal/auth/strategy"
	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/database"
	"github.com/Salam4nder/user/internal/email"
	"github.com/Salam4nder/user/internal/event"
	"github.com/Salam4nder/user/internal/grpc/interceptors"
	"github.com/Salam4nder/user/internal/grpc/server"
	"github.com/Salam4nder/user/internal/observability/metrics"
	"github.com/Salam4nder/user/internal/observability/otel"
	"github.com/Salam4nder/user/internal/token"
	"github.com/Salam4nder/user/pkg/logger"
	"github.com/Salam4nder/user/proto/gen"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stimtech/go-migration/v2"
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

	// TODO(kg):  Move these to config.
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
	migrationFolder      = "/app/db/migrations"
	natsTimeout          = 5 * time.Second
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

	psqlDB, err := sql.Open(cfg.PSQL.Driver(), cfg.PSQL.Addr())
	exitOnError(ctx, err)
	if err = database.HealthCheck(ctx, psqlDB, 5 /*max tries*/); err != nil {
		exitOnError(ctx, err)
	}

	if err = migration.New(
		psqlDB,
		migration.Config{MigrationFolder: migrationFolder},
	).Migrate(); err != nil {
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

	// Worker.
	go event.NewWorker(email.NewNoOpSender()).Work(ctx, natsChan)

	// Token maker.
	tokenMaker, err := token.BootstrapPasetoMaker(
		accessTokenDuration,
		refreshTokenDuration,
		[]byte(cfg.SymmetricKey))
	if err != nil {
		exitOnError(ctx, err)
	}

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
	userServer, err := server.NewUserServer(
		psqlDB,
		healthServer,
		natsClient,
		strategy.NewCredentials(psqlDB, natsClient),
		tokenMaker,
	)
	exitOnError(ctx, err)
	gen.RegisterIdentityServer(grpcServer, userServer)
	reflection.Register(grpcServer)

	go userServer.MonitorHealth(ctx)

	srvErrChan := make(chan error, 1)
	go func() {
		srvErrChan <- grpcServer.Serve(grpcListener)
	}()
	slog.InfoContext(ctx, "main: gRPC server is running", "address", cfg.Server.GRPCAddr())

	if err = metrics.Register(); err != nil {
		exitOnError(ctx, err)
	}
	http.Handle("/metrics", promhttp.Handler())
	promSrv := http.Server{
		Addr:        "0.0.0.0:8090",
		ReadTimeout: time.Second * 10,
	}
	go func() {
		srvErrChan <- promSrv.ListenAndServe()
	}()
	slog.InfoContext(ctx, "main: serving metrics on ", "address", ":8090")

	select {
	case err := <-srvErrChan:
		slog.ErrorContext(ctx, "gRPC server error", "error", err)
	case <-ctx.Done():
		slog.InfoContext(ctx, "main: context done, shutting down...")
	}
	grpcServer.GracefulStop()
	err = errors.Join(err, psqlDB.Close())
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
