package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

const levelTraceOffset slog.Level = 4 // distance from Debug/Error to trace/fatal custom levels

func parseLevel(logLevel string) (slog.Level, error) {
	switch strings.ToLower(logLevel) {
	case "trace":
		return slog.LevelDebug - levelTraceOffset, nil
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	case "fatal", "panic":
		return slog.LevelError + levelTraceOffset, nil
	default:
		return 0, fmt.Errorf("unknown log level %q", logLevel)
	}
}

// InitLog creates and returns a logger configured at the given log level.
func InitLog(logLevel string) (*slog.Logger, error) {
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
	return slog.New(handler), nil
}

// InitLogFromEnv creates and returns a logger configured from the LOG_LEVEL
// environment variable. If the variable is not set, it defaults to info.
func InitLogFromEnv() (*slog.Logger, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		fmt.Fprintln(os.Stderr, "warn: LOG_LEVEL not set, defaulting to info")
		logLevel = "info"
	}
	return InitLog(logLevel)
}
