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

func (s slogDefaultLogger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	s.getSlogLogger().LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
}

func (s slogDefaultLogger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	s.getSlogLogger().LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

func (s slogDefaultLogger) Warn(ctx context.Context, err error, attrs ...slog.Attr) {
	s.getSlogLogger().LogAttrs(ctx, slog.LevelWarn, err.Error(), attrs...)
}

func (s slogDefaultLogger) Error(ctx context.Context, err error, attrs ...slog.Attr) {
	s.getSlogLogger().LogAttrs(ctx, slog.LevelError, err.Error(), attrs...)
}
