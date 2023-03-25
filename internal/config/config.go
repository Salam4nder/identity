package config

import (
	"fmt"

	"github.com/plaid/go-envvar/envvar"
)

// Application is the application configuration.
type Application struct {
	Mongo  MongoDB    `envvar:"MONGO_"`
	Server GRPCServer `envvar:"GRPC_SERVER_"`
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

// MongoDB holds the MongoDB configuration.
type MongoDB struct {
	Host       string `envvar:"HOST" default:"localhost"`
	Port       string `envvar:"PORT" default:"27017"`
	Username   string `envvar:"USERNAME" default:"root"`
	Password   string `envvar:"PASSWORD" default:"root"`
	Name       string `envvar:"NAME" default:"user"`
	Collection string `envvar:"COLLECTION" default:"users"`
}

// GRPCServer holds the gRPC server configuration.
type GRPCServer struct {
	Host string `envvar:"HOST" default:"localhost"`
	Port string `envvar:"PORT" default:"8080"`
}

// URI returns the mongoDB connection string.
func (m *MongoDB) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		m.Username,
		m.Password,
		m.Host,
		m.Port)
}

// Addr returns the gRPC server address.
func (g *GRPCServer) Addr() string {
	return fmt.Sprintf("%s:%s", g.Host, g.Port)
}
