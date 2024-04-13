// Package config provides the application configuration.
// Currently using yaml package to parse environment variables.
// Planning to switch to Viper in the future.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Application is the application configuration.
type Application struct {
	Environment string   `yaml:"environment"`
	PSQL        Postgres `yaml:"postgres"`
	Redis       Redis    `yaml:"redis"`
	Server      Server   `yaml:"server"`
}

func (x Application) String() string {
	return fmt.Sprintf(
		"environment: %s, postgres: %v, redis: %v, server: %v",
		x.Environment,
		x.PSQL,
		x.Redis,
		x.Server,
	)
}

// New returns a new application configuration
// Returns an error if any of the environment variables are missing.
func New() (*Application, error) {
	var cfg Application

	f, err := os.Open("config.yaml")
	if err != nil {
		log.Error().Err(err).Msg("config: opening config file")
		return nil, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		log.Error().Err(err).Msg("config: decoding config file")
		return nil, err
	}

	log.Info().Msgf("config: %s", cfg.String())

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

// Redis holds the Redis configuration.
type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Server holds the gRPC server configuration.
type Server struct {
	GRPCHost             string        `yaml:"host"`
	GRPCPort             string        `yaml:"port"`
	SymmetricKey         string        `yaml:"symmetricKey"`
	AccessTokenDuration  time.Duration `yaml:"accessTokenDuration"`
	RefreshTokenDuration time.Duration `yaml:"refreshTokenDuration"`
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
