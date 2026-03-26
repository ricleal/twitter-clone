package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

func parseLevel(logLevel string) (slog.Level, error) {
	switch strings.ToLower(logLevel) {
	case "trace":
		return slog.LevelDebug - 4, nil
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	case "fatal", "panic":
		return slog.LevelError + 4, nil
	default:
		return 0, fmt.Errorf("unknown log level %q", logLevel)
	}
}

// InitLog initializes the logger with the given log level and sets it as default.
func InitLog(ctx context.Context, logLevel string) (context.Context, error) {
	level, err := parseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s", logLevel)
	}
	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if isatty.IsTerminal(os.Stdout.Fd()) {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))
	return ctx, nil
}

// InitLogFromEnv initializes the logger with the log level from the LOG_LEVEL
// environment variable. If the variable is not set, it defaults to info.
func InitLogFromEnv(ctx context.Context) (context.Context, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		slog.Warn("No environment variable LOG_LEVEL, setting it to info", "level", logLevel)
		logLevel = "info"
	}
	return InitLog(ctx, logLevel)
}
