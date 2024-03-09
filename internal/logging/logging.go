package logging

import (
	"log/slog"
	"os"
)

const (
	// envLocal is the local environment
	envLocal = "local"
)

// New returns a new logger
// If the ENV is local, it will return a logger with a text handler
// If the ENV is anything else, it will return a logger with a JSON handler
func New() *slog.Logger {
	env := os.Getenv("ENV")
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	}

	return log
}
