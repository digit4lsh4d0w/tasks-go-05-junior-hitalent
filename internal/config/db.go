package config

import (
	"errors"
	"fmt"
	"slices"
)

type DBConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

func (c *DBConfig) Validate() error {
	var errs []error

	validDrivers := []string{"sqlite", "postgres"}

	if !slices.Contains(validDrivers, c.Driver) {
		errs = append(errs, &ValidationError{Field: "database.driver", Msg: fmt.Sprintf("unknown driver %q", c.Driver)})
	}

	if c.DSN == "" {
		errs = append(errs, &ValidationError{Field: "database.dsn", Msg: "dsn is required"})
	}

	return errors.Join(errs...)
}
