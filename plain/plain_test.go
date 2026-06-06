package plain_test

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/Siroshun09/logs/v2"
	"github.com/Siroshun09/logs/v2/plain"
)

func Test_plainLogger_Debug(t *testing.T) {
	tests := []struct {
		name        string
		opt         *plain.Option
		call        func(ctx context.Context, logger logs.Logger)
		expectedLog string
	}{
		{
			name: "default option -> not printed",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log")
			},
			expectedLog: "",
		},
		{
			name: "empty option -> not printed",
			opt:  &plain.Option{},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log")
			},
			expectedLog: "",
		},
		{
			name: "debug enabled -> printed",
			opt: &plain.Option{
				Debug: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log")
			},
			expectedLog: "test log\n",
		},
		{
			name: "debug enabled with print level -> printed with level",
			opt: &plain.Option{
				Debug:      true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log")
			},
			expectedLog: "DEBUG: test log\n",
		},
		{
			name: "debug enabled -> log with attrs -> printed without attrs",
			opt: &plain.Option{
				Debug: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log", slog.String("key", "value"))
			},
			expectedLog: "test log\n",
		},
		{
			name: "debug enabled with print attrs -> log with attrs -> printed with attrs",
			opt: &plain.Option{
				Debug:      true,
				PrintAttrs: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log", slog.String("key", "value"))
			},
			expectedLog: "test log key=value\n",
		},
		{
			name: "debug enabled with print attrs and level -> log with attrs -> printed with attrs and level",
			opt: &plain.Option{
				Debug:      true,
				PrintAttrs: true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log", slog.String("key", "value"))
			},
			expectedLog: "DEBUG: test log key=value\n",
		},
		{
			name: "debug enabled with print attrs and level -> log without attrs -> printed with level",
			opt: &plain.Option{
				Debug:      true,
				PrintAttrs: true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log")
			},
			expectedLog: "DEBUG: test log\n",
		},
		{
			name: "debug disabled -> not printed",
			opt: &plain.Option{
				Debug: false,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log")
			},
			expectedLog: "",
		},
		{
			name: "debug disabled -> log with attrs -> not printed",
			opt: &plain.Option{
				Debug: false,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Debug(ctx, "test log", slog.String("key", "value"))
			},
			expectedLog: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.Builder{}
			tt.call(t.Context(), plain.NewPlainLogger(&b, tt.opt))

			if b.String() != tt.expectedLog {
				t.Errorf("plainLogger.Debug() = %v, want %v", b.String(), tt.expectedLog)
			}
		})
	}
}

func Test_plainLogger_Info(t *testing.T) {
	tests := []struct {
		name        string
		opt         *plain.Option
		call        func(ctx context.Context, logger logs.Logger)
		expectedLog string
	}{
		{
			name: "default option -> printed",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Info(ctx, "test log")
			},
			expectedLog: "test log\n",
		},
		{
			name: "empty option -> printed",
			opt:  &plain.Option{},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Info(ctx, "test log")
			},
			expectedLog: "test log\n",
		},
		{
			name: "debug enabled -> printed",
			opt: &plain.Option{
				Debug: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Info(ctx, "test log")
			},
			expectedLog: "test log\n",
		},
		{
			name: "with print level -> printed with level",
			opt: &plain.Option{
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Info(ctx, "test log")
			},
			expectedLog: "INFO: test log\n",
		},
		{
			name: "default option -> log with attrs -> printed without attrs",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Info(ctx, "test log", slog.String("key", "value"))
			},
			expectedLog: "test log\n",
		},
		{
			name: "with print attrs -> log with attrs -> printed with attrs",
			opt: &plain.Option{
				PrintAttrs: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Info(ctx, "test log", slog.String("key", "value"))
			},
			expectedLog: "test log key=value\n",
		},
		{
			name: "with print attrs and level -> log with attrs -> printed with attrs and level",
			opt: &plain.Option{
				PrintAttrs: true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Info(ctx, "test log", slog.String("key", "value"))
			},
			expectedLog: "INFO: test log key=value\n",
		},
		{
			name: "with print attrs and level -> log without attrs -> printed with level",
			opt: &plain.Option{
				PrintAttrs: true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Info(ctx, "test log")
			},
			expectedLog: "INFO: test log\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.Builder{}
			tt.call(t.Context(), plain.NewPlainLogger(&b, tt.opt))

			if b.String() != tt.expectedLog {
				t.Errorf("plainLogger.Info() = %v, want %v", b.String(), tt.expectedLog)
			}
		})
	}
}

func Test_plainLogger_Warn(t *testing.T) {
	err := errors.New("test error")
	tests := []struct {
		name        string
		opt         *plain.Option
		call        func(ctx context.Context, logger logs.Logger)
		expectedLog string
	}{
		{
			name: "default option -> printed",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, err)
			},
			expectedLog: "test error\n",
		},
		{
			name: "nil error -> printed as <nil>",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, nil)
			},
			expectedLog: "<nil>\n",
		},
		{
			name: "empty option -> printed",
			opt:  &plain.Option{},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, err)
			},
			expectedLog: "test error\n",
		},
		{
			name: "debug enabled -> printed",
			opt: &plain.Option{
				Debug: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, err)
			},
			expectedLog: "test error\n",
		},
		{
			name: "with print level -> printed with level",
			opt: &plain.Option{
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, err)
			},
			expectedLog: "WARN: test error\n",
		},
		{
			name: "default option -> log with attrs -> printed without attrs",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, err, slog.String("key", "value"))
			},
			expectedLog: "test error\n",
		},
		{
			name: "with print attrs -> log with attrs -> printed with attrs",
			opt: &plain.Option{
				PrintAttrs: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, err, slog.String("key", "value"))
			},
			expectedLog: "test error key=value\n",
		},
		{
			name: "with print attrs and level -> log with attrs -> printed with attrs and level",
			opt: &plain.Option{
				PrintAttrs: true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, err, slog.String("key", "value"))
			},
			expectedLog: "WARN: test error key=value\n",
		},
		{
			name: "with print attrs and level -> log without attrs -> printed with level",
			opt: &plain.Option{
				PrintAttrs: true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Warn(ctx, err)
			},
			expectedLog: "WARN: test error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.Builder{}
			tt.call(t.Context(), plain.NewPlainLogger(&b, tt.opt))

			if b.String() != tt.expectedLog {
				t.Errorf("plainLogger.Warn() = %v, want %v", b.String(), tt.expectedLog)
			}
		})
	}
}

func Test_plainLogger_Error(t *testing.T) {
	err := errors.New("test error")
	tests := []struct {
		name        string
		opt         *plain.Option
		call        func(ctx context.Context, logger logs.Logger)
		expectedLog string
	}{
		{
			name: "default option -> printed",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, err)
			},
			expectedLog: "test error\n",
		},
		{
			name: "nil error -> printed as <nil>",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, nil)
			},
			expectedLog: "<nil>\n",
		},
		{
			name: "empty option -> printed",
			opt:  &plain.Option{},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, err)
			},
			expectedLog: "test error\n",
		},
		{
			name: "debug enabled -> printed",
			opt: &plain.Option{
				Debug: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, err)
			},
			expectedLog: "test error\n",
		},
		{
			name: "with print level -> printed with level",
			opt: &plain.Option{
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, err)
			},
			expectedLog: "ERROR: test error\n",
		},
		{
			name: "default option -> log with attrs -> printed without attrs",
			opt:  nil,
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, err, slog.String("key", "value"))
			},
			expectedLog: "test error\n",
		},
		{
			name: "with print attrs -> log with attrs -> printed with attrs",
			opt: &plain.Option{
				PrintAttrs: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, err, slog.String("key", "value"))
			},
			expectedLog: "test error key=value\n",
		},
		{
			name: "with print attrs and level -> log with attrs -> printed with attrs and level",
			opt: &plain.Option{
				PrintAttrs: true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, err, slog.String("key", "value"))
			},
			expectedLog: "ERROR: test error key=value\n",
		},
		{
			name: "with print attrs and level -> log without attrs -> printed with level",
			opt: &plain.Option{
				PrintAttrs: true,
				PrintLevel: true,
			},
			call: func(ctx context.Context, logger logs.Logger) {
				logger.Error(ctx, err)
			},
			expectedLog: "ERROR: test error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.Builder{}
			tt.call(t.Context(), plain.NewPlainLogger(&b, tt.opt))

			if b.String() != tt.expectedLog {
				t.Errorf("plainLogger.Error() = %v, want %v", b.String(), tt.expectedLog)
			}
		})
	}
}
