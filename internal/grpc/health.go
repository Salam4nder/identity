package grpc

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"
)

// MonitorHealth will ping all [UserServer] dependencies every 5 seconds and update the health status.
func (x *UserServer) MonitorHealth(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			x.health.Shutdown()
			return
		case <-time.After(5 * time.Second):
			var unhealthy bool
			if err := x.storage.PingContext(ctx, 1); err != nil {
				unhealthy = true
				// Only log if the context is not done.
				if ctx.Err() == nil {
					slog.Warn("health: pinging storage failed, setting health status to not serving")
				}
			}
			if !x.natsConn.IsConnected() {
				unhealthy = true
				if ctx.Err() == nil {
					slog.Warn("health: nats is not connected, setting health status to not serving")
				}
			}
			if unhealthy {
				// Empty string means the whole service.
				x.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			} else {
				x.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
			}
		}
	}
}
