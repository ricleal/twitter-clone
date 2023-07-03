package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLog(ctx context.Context, logLevel string) (context.Context, error) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s", logLevel)
	}
	logger := zerolog.New(os.Stdout)
	logger = logger.With().Timestamp().Logger().Level(level)
	if isatty.IsTerminal(os.Stdout.Fd()) {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	ctx = logger.WithContext(ctx)
	log.Logger = logger
	return ctx, nil
}

// InitLogFromEnv initializes the logger with the log level from the LOG_LEVEL
// environment variable. If the variable is not set, it defaults to info.
func InitLogFromEnv(ctx context.Context) (context.Context, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		log.Warn().Str("level", logLevel).Msg("No environment variable LOG_LEVEL, setting it to info")
		logLevel = "info"
	}
	return InitLog(ctx, logLevel)
}
