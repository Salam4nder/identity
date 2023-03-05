package config

import (
	"fmt"

	"github.com/plaid/go-envvar/envvar"
)

// Application is the application configuration
type Application struct {
	Mongo MongoDB
}

// New returns a new application configuration
// Returns an error if any of the environment variables are missing
func New() (*Application, error) {
	var cfg Application
	if err := envvar.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// MongoDB holds the MongoDB configuration
type MongoDB struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Name       string
	Collection string
}

// URI returns the mongoDB connection string
func (m *MongoDB) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		m.Username,
		m.Password,
		m.Host,
		m.Port)
}
