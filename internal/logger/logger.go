package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/s-588/tms/internal/config"
)

// SetupSLog() is a function that configures default slog logger with given cfg.
// It returns os.File.Close() function that must be closed when app exits.
// It handles config.LoggerConfig.Level by itself, default level is slog.LevelInfo.
// If logger level is "DEBUG" slogs will include source.
func SetupSLog(cfg config.LoggerConfig) (func() error, error) {
	opts := &slog.HandlerOptions{
		Level:     parseLogLevel(cfg.Level),
		AddSource: cfg.Level == "DEBUG",
	}

	h := slog.NewTextHandler(os.Stdout, opts)
	if cfg.File != "" {
		path := filepath.Clean(cfg.File)

		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return nil, fmt.Errorf("can't setup slog logger: %w", err)
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("can't setup slog logger: %w", err)
		}

		h = slog.NewTextHandler(io.MultiWriter(os.Stdout, f), opts)
		slog.SetDefault(slog.New(h))
		return f.Close, nil
	}

	slog.SetDefault(slog.New(h))
	return func() error { return nil }, nil
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "ERROR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
