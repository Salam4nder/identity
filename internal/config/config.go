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
func New() (*Application, error) {
	var cfg Application
	if err := envvar.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// MongoDB holds the MongoDB configuration
type MongoDB struct {
	Host     string
	Port     string
	Username string
	Password string
}

// URI returns the database connection string
func (d *MongoDB) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port)
}
