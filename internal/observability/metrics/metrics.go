package metrics

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	UsersRegistered = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "user",
		Subsystem: "api",
		Name:      "users_registered_total",
		Help:      "Number of registered users - increases on user creation",
	})

	UsersActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "api",
		Name:      "users_active",
		Help:      "Number of active users - for now increases on user creation and decreases on deletion",
	})
)

// Register will register all collectors defined in metrics.go.
func Register() error {
	collectors := []prometheus.Collector{
		UsersRegistered,
		UsersActive,
	}
	var errs []error
	for i := range collectors {
		if err := prometheus.Register(collectors[i]); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
