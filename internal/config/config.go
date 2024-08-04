package config

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// Application is the application configuration.
type Application struct {
	Environment  string `yaml:"environment"`
	SymmetricKey string `yaml:"symmetricKey"`
	// AccessDuration  time.Duration `yaml:"accessDuration"`
	// RefreshDuration time.Duration `yaml:"refreshDuration"`
	PSQL   Postgres `yaml:"postgres"`
	NATS   NATS     `yaml:"nats"`
	Server Server   `yaml:"server"`
}

// New returns a new application configuration
// Returns an error if any of the environment variables are missing.
func New() (*Application, error) {
	var cfg Application

	f, err := os.Open("config.yaml")
	if err != nil {
		slog.Error("config: opening config file", "err", err)
		return nil, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		slog.Error("config: decoding config file", "err", err)
		return nil, err
	}

	return &cfg, nil
}

// Postgres holds the Postgres configuration.
type Postgres struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	Name            string `yaml:"db"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	ApplicationName string `yaml:"applicationName"`
}

// NATS holds the NATS configuration.
type NATS struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Server holds the gRPC server configuration.
type Server struct {
	GRPCHost string `yaml:"host"`
	GRPCPort string `yaml:"port"`
}

// Addr returns the PSQL connection string.
func (x *Postgres) Addr() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&application_name=%s",
		x.User,
		x.Password,
		x.Host,
		x.Port,
		x.Name,
		x.ApplicationName,
	)
}

// Driver returns the database driver name.
func (x *Postgres) Driver() string {
	return "postgres"
}

// Addr returns the NATS connection string.
func (x NATS) Addr() string {
	return fmt.Sprintf("nats://%s:%s", x.Host, x.Port)
}

// GRPCAddr returns the gRPC server address.
func (x *Server) GRPCAddr() string {
	return fmt.Sprintf("%s:%s", x.GRPCHost, x.GRPCPort)
}

// PSQLTestConfig is used to connect to the unit test db.
func PSQLTestConfig() Postgres {
	return Postgres{
		Host:     "localhost",
		Port:     "54321",
		Name:     "unit-test-user-db",
		User:     "test",
		Password: "unit-test-pw",
	}
}
