package logs

import (
	"context"
	"log/slog"
)

// NewLoggerWithSlog creates Logger using slog.Logger.
func NewLoggerWithSlog(logger *slog.Logger) Logger {
	return slogDefaultLogger{delegate: logger}
}

type slogDefaultLogger struct {
	delegate *slog.Logger
}

func (s slogDefaultLogger) getSlogLogger() *slog.Logger {
	if s.delegate != nil {
		return s.delegate
	}
	return slog.Default()
}

func (s slogDefaultLogger) Debug(ctx context.Context, msg string) {
	s.getSlogLogger().DebugContext(ctx, msg)
}

func (s slogDefaultLogger) Info(ctx context.Context, msg string) {
	s.getSlogLogger().InfoContext(ctx, msg)
}

func (s slogDefaultLogger) Warn(ctx context.Context, err error) {
	s.getSlogLogger().WarnContext(ctx, err.Error())
}

func (s slogDefaultLogger) Error(ctx context.Context, err error) {
	s.getSlogLogger().ErrorContext(ctx, err.Error())
}
