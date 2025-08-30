package logs

import "context"

// Logger includes functions to log messages or errors.
type Logger interface {
	// Debug logs a message as debug level.
	Debug(ctx context.Context, msg string)
	// Info logs a message as info level.
	Info(ctx context.Context, msg string)
	// Warn logs an error as warn level.
	Warn(ctx context.Context, err error)
	// Warnf logs an error as warn level.
	Warnf(ctx context.Context, format string, args ...any)
	// Error logs an error as error level.
	Error(ctx context.Context, err error)
	// Errorf logs an error as error level.
	Errorf(ctx context.Context, format string, args ...any)
}
