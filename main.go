package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/db"
	internalGRPC "github.com/Salam4nder/user/internal/grpc"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/internal/task"
	grpcUtil "github.com/Salam4nder/user/pkg/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stimtech/go-migration"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
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

	otelShutdown, err := setupOTelSDK(ctx)

	db, err := db.NewSQLDatabase(ctx, cfg.PSQL)
	exitWithError(err)
	migration := migration.New(db.DB(), zap.NewNop())
	if err = migration.Migrate(); err != nil {
		exitWithError(err)
	}

	// TODO: config this better.
	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.Redis.Addr(),
	}
	taskCreator := task.NewRedisTaskCreator(redisOpt)
	taskProcessor := task.NewRedisTaskProcessor(db, redisOpt)
	go func() {
		if err := taskProcessor.Process(); err != nil {
			exitWithError(err)
		}
	}()

	userServer, err := internalGRPC.NewUserServer(db, taskCreator, cfg.Server)
	exitWithError(err)
	grpcListener, err := net.Listen("tcp", cfg.Server.GRPCAddr())
	exitWithError(err)
	grpcServer := grpc.NewServer(
		// grpc.StatsHandler(otelgrpc.NewServerHandler()),
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("signal received, shutting down...")
	grpcServer.GracefulStop()
	if err := otelShutdown(ctx); err != nil {
		exitWithError(err)
	}
	log.Info().Msg("service gracefully stopped")
}

func exitWithError(err error) {
	if err != nil {
		log.Error().Err(err).Msg("main: exit with error")
		os.Exit(1)
	}
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider() (*trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	grpcExporter, err := otlptracegrpc.New(context.Background())

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithBatcher(grpcExporter),
	)
	return traceProvider, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)
	return meterProvider, nil
}
