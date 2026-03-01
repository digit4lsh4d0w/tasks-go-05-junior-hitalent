package config

import (
	"errors"
	"fmt"
	"slices"
)

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

func (c *LogConfig) Validate() error {
	var errs []error

	validLevels := []string{"debug", "info", "warning", "error"}
	validOutputs := []string{"stdout", "file", "both"}
	validFormats := []string{"text", "json"}

	if !slices.Contains(validLevels, c.Level) {
		errs = append(errs, &ValidationError{Field: "log.level", Msg: fmt.Sprintf("unknown value %q", c.Level)})
	}

	if !slices.Contains(validOutputs, c.Output) {
		errs = append(errs, &ValidationError{Field: "log.output", Msg: fmt.Sprintf("unknown value %q", c.Output)})
	}

	if !slices.Contains(validFormats, c.Format) {
		errs = append(errs, &ValidationError{Field: "log.format", Msg: fmt.Sprintf("unknown value %q", c.Format)})
	}

	if (c.Output == "file" || c.Output == "both") && c.Path == "" {
		errs = append(errs, &ValidationError{Field: "log.path", Msg: "path is required when output is \"file\" or \"both\""})
	}

	return errors.Join(errs...)
}
