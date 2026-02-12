package logs

import (
	"context"
	"log/slog"
)

// Logger includes functions to log messages or errors.
type Logger interface {
	// Debug logs a message as debug level.
	Debug(ctx context.Context, msg string, attrs ...slog.Attr)
	// Info logs a message as info level.
	Info(ctx context.Context, msg string, attrs ...slog.Attr)
	// Warn logs an error as warn level.
	Warn(ctx context.Context, err error, attrs ...slog.Attr)
	// Error logs an error as error level.
	Error(ctx context.Context, err error, attrs ...slog.Attr)
}
