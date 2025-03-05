package logging

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func NewLogger(logLevel string) *slog.Logger {

	var slogLevel slog.Level
	switch logLevel {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warning":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	}

	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slogLevel},
	)
	Logger = slog.New(handler)
	return Logger
}
