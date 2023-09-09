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
	Server      Server      `envvar:"SERVER_"`
	Service     UserService `envvar:"USER_SERVICE_"`
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

// Server holds the gRPC server configuration.
type Server struct {
	GRPCHost string `envvar:"GRPC_HOST" default:"localhost"`
	GRPCPort string `envvar:"GRPC_PORT" default:"8080"`
	HTTPHost string `envvar:"HTTP_HOST" default:"localhost"`
	HTTPPort string `envvar:"HTTP_PORT" default:"8081"`
}

// URI returns the mongoDB connection string.
func (dbCfg *Postgres) URI() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&application_name=%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
		dbCfg.ApplicationName,
	)
}

// Driver returns the database driver name.
func (dbCfg *Postgres) Driver() string {
	return "postgres"
}

// GRPCAddr returns the gRPC server address.
func (srvCfg *Server) GRPCAddr() string {
	return fmt.Sprintf("%s:%s", srvCfg.GRPCHost, srvCfg.GRPCPort)
}

// HTTPAddr returns the gRPC gateway server address.
func (srvCfg *Server) HTTPAddr() string {
	return fmt.Sprintf("%s:%s", srvCfg.HTTPHost, srvCfg.HTTPPort)
}
