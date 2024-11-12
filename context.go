package logs

import "context"

type contextKey uint8

const (
	loggerKey contextKey = iota + 1
)

// WithContext puts the Logger into context.Context.
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext takes the Logger from context.Context.
//
// If the given context.Context does not have the Logger, this function returns Default.
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	return Default()
}

// Debug logs a message as debug level using Logger obtained from FromContext.
func Debug(ctx context.Context, msg string) {
	FromContext(ctx).Debug(ctx, msg)
}

// Info logs a message as info level using Logger obtained from FromContext.
func Info(ctx context.Context, msg string) {
	FromContext(ctx).Info(ctx, msg)
}

// Warn logs an error as warn level using Logger obtained from FromContext.
func Warn(ctx context.Context, err error) {
	FromContext(ctx).Warn(ctx, err)
}

// Error logs an error as error level using Logger obtained from FromContext.
func Error(ctx context.Context, err error) {
	FromContext(ctx).Error(ctx, err)
}
