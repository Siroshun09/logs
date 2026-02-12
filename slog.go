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
	s.getSlogLogger().DebugContext(ctx, msg, s.toAnyArray(attrs)...)
}

func (s slogDefaultLogger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	s.getSlogLogger().InfoContext(ctx, msg, s.toAnyArray(attrs)...)
}

func (s slogDefaultLogger) Warn(ctx context.Context, err error, attrs ...slog.Attr) {
	s.getSlogLogger().WarnContext(ctx, err.Error(), s.toAnyArray(attrs)...)
}

func (s slogDefaultLogger) Error(ctx context.Context, err error, attrs ...slog.Attr) {
	s.getSlogLogger().ErrorContext(ctx, err.Error(), s.toAnyArray(attrs)...)
}

func (s slogDefaultLogger) toAnyArray(attrs []slog.Attr) []any {
	if attrs == nil {
		return nil
	}

	anyAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		anyAttrs[i] = attr
	}

	return anyAttrs
}
