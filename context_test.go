package logs

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"
)

func TestWithContext(t *testing.T) {
	ctx := WithContext(context.Background(), Default())
	if _, ok := ctx.Value(loggerKey).(Logger); !ok {
		t.Errorf("logger not found in context")
	}
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
			ctx := tt.ctxFunc(context.Background())
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
	log any
}

func (r recordingLogger) Debug(_ context.Context, msg string) {
	r.lastLog.log = debugLog(msg)
}

func (r recordingLogger) Info(_ context.Context, msg string) {
	r.lastLog.log = infoLog(msg)
}

func (r recordingLogger) Warn(_ context.Context, err error) {
	r.lastLog.log = warnLog(err)
}

func (r recordingLogger) Error(_ context.Context, err error) {
	r.lastLog.log = errorLog(err)
}

func TestLog(t *testing.T) {
	lastLog := lastLogHolder{}
	ctx := WithContext(context.Background(), &recordingLogger{lastLog: &lastLog})
	warn := errors.New("warn")
	err := errors.New("error")

	tests := []struct {
		name        string
		f           func(ctx context.Context)
		expectedLog any
	}{
		{
			name: "success: debug",
			f: func(ctx context.Context) {
				Debug(ctx, "debug")
			},
			expectedLog: debugLog("debug"),
		},
		{
			name: "success: info",
			f: func(ctx context.Context) {
				Info(ctx, "info")
			},
			expectedLog: infoLog("info"),
		},
		{
			name: "success: warn",
			f: func(ctx context.Context) {
				Warn(ctx, warn)
			},
			expectedLog: warnLog(warn),
		},
		{
			name: "success: error",
			f: func(ctx context.Context) {
				Error(ctx, err)
			},
			expectedLog: errorLog(err),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f(ctx)
			if !reflect.DeepEqual(tt.expectedLog, lastLog.log) {
				t.Errorf("lastLogHolder = %v, want %v", lastLog.log, tt.expectedLog)
			}
		})
	}
}
