package errorlogs

import (
	"context"
	"log/slog"
	"slices"

	"github.com/Siroshun09/logs/v2"
	"github.com/Siroshun09/serrors/v2"
)

const DefaultStackTraceAttrKey = "stack_trace"

// NewLogger creates a new logs.Logger.
func NewLogger(out logs.Logger, option *Option) logs.Logger {
	opt := Option{}
	if option != nil {
		opt = *option
	}

	return &logger{
		dedicated: out,
		opt:       opt,
	}
}

// Option is the option for logger implementation.
type Option struct {
	// PrintStackTraceOnWarn is whether to print stack trace on Warn.
	PrintStackTraceOnWarn bool
	// PrintCurrentStackTraceIfNotAttached is whether to print the current stack trace if the error does not have a stack trace.
	PrintCurrentStackTraceIfNotAttached bool
	// StackTraceAttrKey is the key of the stack trace attribute.
	StackTraceAttrKey string
	// ErrorAttrsFunc is the function to modify the attributes of the error.
	ErrorAttrsFunc func(err error, attrs []slog.Attr) []slog.Attr
}

type logger struct {
	dedicated logs.Logger
	opt       Option
}

func (l *logger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l == nil {
		return
	}

	l.dedicated.Debug(ctx, msg, attrs...)
}

func (l *logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	if l == nil {
		return
	}

	l.dedicated.Info(ctx, msg, attrs...)
}

func (l *logger) Warn(ctx context.Context, err error, attrs ...slog.Attr) {
	if l == nil {
		return
	}

	attrs = l.appendAttrs(attrs, err, l.opt.PrintStackTraceOnWarn)
	l.dedicated.Warn(ctx, err, attrs...)
}

func (l *logger) Error(ctx context.Context, err error, attrs ...slog.Attr) {
	if l == nil {
		return
	}

	attrs = l.appendAttrs(attrs, err, true)
	l.dedicated.Error(ctx, err, attrs...)
}

func (l *logger) appendAttrs(attrs []slog.Attr, err error, includeStackTrace bool) []slog.Attr {
	// Clip so that appending never mutates the caller's backing array.
	ret := slices.Clip(attrs)

	if includeStackTrace {
		key := l.opt.StackTraceAttrKey
		if key == "" {
			key = DefaultStackTraceAttrKey
		}

		if attached, exists := serrors.GetAttachedStackTrace(err); exists {
			ret = append(ret, slog.Any(key, attached))
		} else if l.opt.PrintCurrentStackTraceIfNotAttached {
			ret = append(ret, slog.Any(key, serrors.GetCurrentStackTrace()))
		}
	}

	var errAttrs []slog.Attr
	for _, attr := range serrors.GetAttrs(err) {
		errAttrs = append(errAttrs, attr)
	}

	if l.opt.ErrorAttrsFunc != nil {
		errAttrs = l.opt.ErrorAttrsFunc(err, errAttrs)
	}

	return append(ret, errAttrs...)
}
