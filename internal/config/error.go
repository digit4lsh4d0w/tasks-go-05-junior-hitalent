package config

import "errors"

var (
	ErrConfigNotFound = errors.New("config file not found")
	ErrConfigRead     = errors.New("failed to read config file")
	ErrConfigParse    = errors.New("failed to parse config file")
)

type ConfigError struct {
	Op   string
	Path string
	Err  error
}

func (e *ConfigError) Error() string {
	return e.Err.Error()
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}
