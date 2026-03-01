package slog

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"task-5/internal/config"
	"task-5/internal/logger"
)

var _ logger.Logger = (*Logger)(nil)

type Logger struct {
	*slog.Logger
	closer io.Closer
}

func toSlogLevel(l string) slog.Leveler {
	switch strings.ToLower(l) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func openLogFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log file %q: %w", path, err)
	}

	return f, nil
}

func newHandler(format string, w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	switch strings.ToLower(format) {
	case "json":
		return slog.NewJSONHandler(w, opts)
	default:
		return slog.NewTextHandler(w, opts)
	}
}

func New(cfg *config.LogConfig) (*Logger, error) {
	opts := &slog.HandlerOptions{
		// AddSource: true,
		AddSource: cfg.AddSource,
		Level:     toSlogLevel(cfg.Level),
	}

	var handler slog.Handler
	var closer io.Closer

	switch cfg.Output {
	case "file":
		f, err := openLogFile(cfg.Path)
		if err != nil {
			return nil, fmt.Errorf("logger: cannot open log file %q: %w", cfg.Path, err)
		}

		handler = newHandler(cfg.Format, f, opts)
		closer = f
	case "both":
		f, err := openLogFile(cfg.Path)
		if err != nil {
			return nil, fmt.Errorf("logger: cannot open log file %q: %w", cfg.Path, err)
		}

		handler = slog.NewMultiHandler(
			newHandler(cfg.Format, os.Stdout, opts),
			newHandler(cfg.Format, f, opts),
		)
		closer = f
	default:
		handler = newHandler(cfg.Format, os.Stdout, opts)
	}

	return &Logger{Logger: slog.New(handler), closer: closer}, nil
}

func (l *Logger) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

func (l *Logger) With(args ...any) logger.Logger {
	return &Logger{Logger: l.Logger.With(args...)}
}
