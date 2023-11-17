// Package config provides the application configuration.
// Currently using envvar package to parse environment variables.
// Planning to switch to Viper in the future.
package config

import (
	"fmt"
	"time"

	"github.com/plaid/go-envvar/envvar"
)

// Application is the application configuration.
type Application struct {
	Environment string      `envvar:"ENVIRONMENT"`
	PSQL        Postgres    `envvar:"POSTGRES_"`
	Redis       Redis       `envvar:"REDIS_"`
	Server      Server      `envvar:"SERVER_"`
	UserService UserService `envvar:"USER_SERVICE_"`
}

// New returns a new application configuration
// Returns an error if any of the environment variables are missing.
func New() (*Application, error) {
	var cfg Application

	if err := envvar.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// UserService holds the user service configuration.
type UserService struct {
	SymmetricKey         string        `envvar:"SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `envvar:"ACCESS_TOKEN_DURATION" default:"1h"`
	RefreshTokenDuration time.Duration `envvar:"REFRESH_TOKEN_DURATION" default:"24h"`
}

// Postgres holds the Postgres configuration.
type Postgres struct {
	Host            string `envvar:"HOST" default:"postgres"`
	Port            string `envvar:"PORT" default:"5432"`
	Name            string `envvar:"DB" default:"user"`
	User            string `envvar:"USER" default:"admin"`
	Password        string `envvar:"PASSWORD" default:"password"`
	ApplicationName string `envvar:"APPLICATION_NAME" default:"user"`
}

// Redis holds the Redis configuration.
type Redis struct {
	Host string `envvar:"HOST" default:"0.0.0.0"`
	Port string `envvar:"PORT" default:"6379"`
}

// Server holds the gRPC server configuration.
type Server struct {
	GRPCHost string `envvar:"GRPC_HOST" default:"localhost"`
	GRPCPort string `envvar:"GRPC_PORT" default:"8080"`
	HTTPHost string `envvar:"HTTP_HOST" default:"localhost"`
	HTTPPort string `envvar:"HTTP_PORT" default:"8081"`
}

// Addr returns the mongoDB connection string.
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

// Addr returns the Redis connection string.
func (x *Redis) Addr() string {
	return fmt.Sprintf("%s:%s", x.Host, x.Port)
}

// GRPCAddr returns the gRPC server address.
func (x *Server) GRPCAddr() string {
	return fmt.Sprintf("%s:%s", x.GRPCHost, x.GRPCPort)
}

// HTTPAddr returns the gRPC gateway server address.
func (x *Server) HTTPAddr() string {
	return fmt.Sprintf("%s:%s", x.HTTPHost, x.HTTPPort)
}
