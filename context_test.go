package logs

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"
)

func TestWithContext(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := WithContext(t.Context(), Default())
		if _, ok := ctx.Value(loggerKey).(Logger); !ok {
			t.Errorf("logger not found in context")
		}
	})

	t.Run("nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		WithContext(t.Context(), nil)
	})
}

func TestFromContext(t *testing.T) {
	logger := NewLoggerWithSlog(&slog.Logger{})
	tests := []struct {
		name    string
		ctxFunc func(ctx context.Context) context.Context
		want    Logger
	}{
		{
			name: "success: not set -> Default",
			ctxFunc: func(ctx context.Context) context.Context {
				return ctx
			},
			want: Default(),
		},

		{
			name: "success: set -> logger",
			ctxFunc: func(ctx context.Context) context.Context {
				return WithContext(ctx, logger)
			},
			want: logger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctxFunc(t.Context())
			if got := FromContext(ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

type debugLog string
type infoLog string
type warnLog error
type errorLog error

type recordingLogger struct {
	lastLog *lastLogHolder
}

type lastLogHolder struct {
	log   any
	attrs []slog.Attr
}

func (r recordingLogger) Debug(_ context.Context, msg string, attrs ...slog.Attr) {
	r.lastLog.log = debugLog(msg)
	r.lastLog.attrs = attrs
}

func (r recordingLogger) Info(_ context.Context, msg string, attrs ...slog.Attr) {
	r.lastLog.log = infoLog(msg)
	r.lastLog.attrs = attrs
}

func (r recordingLogger) Warn(_ context.Context, err error, attrs ...slog.Attr) {
	r.lastLog.log = warnLog(err)
	r.lastLog.attrs = attrs
}

func (r recordingLogger) Error(_ context.Context, err error, attrs ...slog.Attr) {
	r.lastLog.log = errorLog(err)
	r.lastLog.attrs = attrs
}

func TestLog(t *testing.T) {
	lastLog := lastLogHolder{}
	ctx := WithContext(t.Context(), &recordingLogger{lastLog: &lastLog})
	warn := errors.New("warn")
	err := errors.New("error")

	tests := []struct {
		name          string
		f             func(ctx context.Context)
		expectedLog   any
		expectedAttrs []slog.Attr
	}{
		{
			name: "success: debug",
			f: func(ctx context.Context) {
				Debug(ctx, "debug")
			},
			expectedLog:   debugLog("debug"),
			expectedAttrs: nil,
		},
		{
			name: "success: debug with attrs",
			f: func(ctx context.Context) {
				Debug(ctx, "debug", slog.String("key", "value"))
			},
			expectedLog: debugLog("debug"),
			expectedAttrs: []slog.Attr{
				slog.String("key", "value"),
			},
		},
		{
			name: "success: info",
			f: func(ctx context.Context) {
				Info(ctx, "info")
			},
			expectedLog:   infoLog("info"),
			expectedAttrs: nil,
		},
		{
			name: "success: info with attrs",
			f: func(ctx context.Context) {
				Info(ctx, "info", slog.String("key", "value"))
			},
			expectedLog: infoLog("info"),
			expectedAttrs: []slog.Attr{
				slog.String("key", "value"),
			},
		},
		{
			name: "success: warn",
			f: func(ctx context.Context) {
				Warn(ctx, warn)
			},
			expectedLog:   warnLog(warn),
			expectedAttrs: nil,
		},
		{
			name: "success: warn with attrs",
			f: func(ctx context.Context) {
				Warn(ctx, warn, slog.String("key", "value"))
			},
			expectedLog: warnLog(warn),
			expectedAttrs: []slog.Attr{
				slog.String("key", "value"),
			},
		},
		{
			name: "success: error",
			f: func(ctx context.Context) {
				Error(ctx, err)
			},
			expectedLog:   errorLog(err),
			expectedAttrs: nil,
		},
		{
			name: "success: error with attrs",
			f: func(ctx context.Context) {
				Error(ctx, err, slog.String("key", "value"))
			},
			expectedLog: errorLog(err),
			expectedAttrs: []slog.Attr{
				slog.String("key", "value"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f(ctx)
			if !reflect.DeepEqual(tt.expectedLog, lastLog.log) {
				t.Errorf("lastLogHolder = %v, want %v", lastLog.log, tt.expectedLog)
			}
			if !reflect.DeepEqual(tt.expectedAttrs, lastLog.attrs) {
				t.Errorf("lastLogHolder = %v, want %v", lastLog.attrs, tt.expectedAttrs)
			}
		})
	}
}
