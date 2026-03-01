package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogConfig LogConfig      `yaml:"log"`
	DBConfig  DatabaseConfig `yaml:"database"`
}

type LogConfig struct {
	// Variants:
	//  - "debug"
	//  - "info"
	//  - "warning"
	//  - "error"
	Level string `yaml:"level"`

	// Variants:
	//  - "stdout"
	//  - "file"
	//  - "both"
	Output string `yaml:"output"`

	// Variants:
	//  - "text"
	//  - "json"
	Format    string `yaml:"format"`
	Path      string `yaml:"path"`
	AddSource bool   `yaml:"add_source"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrConfigNotFound, path)
		}
		return nil, &ConfigError{Op: "read", Path: path, Err: err}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, &ConfigError{Op: "parse", Path: path, Err: err}
	}

	return &cfg, nil
}
