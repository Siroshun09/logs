package logs

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"
)

func Test_slogDefaultLogger_getSlogLogger(t *testing.T) {
	logger := slog.Logger{}

	tests := []struct {
		name     string
		delegate *slog.Logger
		want     *slog.Logger
	}{
		{
			name:     "success: specified logger",
			delegate: &logger,
			want:     &logger,
		},
		{
			name:     "success: slog.Default",
			delegate: nil,
			want:     slog.Default(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := slogDefaultLogger{delegate: tt.delegate}
			if got := s.getSlogLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSlogLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

type slogHandler struct {
	lastRecord slog.Record
}

func (s *slogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (s *slogHandler) Handle(_ context.Context, record slog.Record) error {
	s.lastRecord = record
	return nil
}

func (s *slogHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return nil
}

func (s *slogHandler) WithGroup(_ string) slog.Handler {
	return nil
}

func Test_slogDefaultLogger_Log(t *testing.T) {
	handler := slogHandler{}
	logger := slogDefaultLogger{slog.New(&handler)}

	tests := []struct {
		name           string
		f              func(ctx context.Context)
		expectedRecord slog.Record
		expectedAttrs  []slog.Attr
	}{
		{
			name: "success: debug",
			f: func(ctx context.Context) {
				logger.Debug(ctx, "debug")
			},
			expectedRecord: slog.Record{
				Message: "debug",
				Level:   slog.LevelDebug,
			},
		},
		{
			name: "success: debug with attrs",
			f: func(ctx context.Context) {
				logger.Debug(ctx, "debug", slog.String("key", "value"))
			},
			expectedRecord: slog.Record{
				Message: "debug",
				Level:   slog.LevelDebug,
			},
			expectedAttrs: []slog.Attr{
				slog.String("key", "value"),
			},
		},
		{
			name: "success: info",
			f: func(ctx context.Context) {
				logger.Info(ctx, "info")
			},
			expectedRecord: slog.Record{
				Message: "info",
				Level:   slog.LevelInfo,
			},
		},
		{
			name: "success: info with attrs",
			f: func(ctx context.Context) {
				logger.Info(ctx, "info", slog.String("key", "value"))
			},
			expectedRecord: slog.Record{
				Message: "info",
				Level:   slog.LevelInfo,
			},
			expectedAttrs: []slog.Attr{
				slog.String("key", "value"),
			},
		},
		{
			name: "success: warn",
			f: func(ctx context.Context) {
				logger.Warn(ctx, errors.New("warn"))
			},
			expectedRecord: slog.Record{
				Message: "warn",
				Level:   slog.LevelWarn,
			},
		},
		{
			name: "success: warn with attrs",
			f: func(ctx context.Context) {
				logger.Warn(ctx, errors.New("warn"), slog.String("key", "value"))
			},
			expectedRecord: slog.Record{
				Message: "warn",
				Level:   slog.LevelWarn,
			},
			expectedAttrs: []slog.Attr{
				slog.String("key", "value"),
			},
		},
		{
			name: "success: warn with nil error",
			f: func(ctx context.Context) {
				logger.Warn(ctx, nil)
			},
			expectedRecord: slog.Record{
				Message: "<nil>",
				Level:   slog.LevelWarn,
			},
		},
		{
			name: "success: error",
			f: func(ctx context.Context) {
				logger.Error(ctx, errors.New("error"))
			},
			expectedRecord: slog.Record{
				Message: "error",
				Level:   slog.LevelError,
			},
		},
		{
			name: "success: error with nil error",
			f: func(ctx context.Context) {
				logger.Error(ctx, nil)
			},
			expectedRecord: slog.Record{
				Message: "<nil>",
				Level:   slog.LevelError,
			},
		},
		{
			name: "success: error with attrs",
			f: func(ctx context.Context) {
				logger.Error(ctx, errors.New("error"), slog.String("key", "value"))
			},
			expectedRecord: slog.Record{
				Message: "error",
				Level:   slog.LevelError,
			},
			expectedAttrs: []slog.Attr{
				slog.String("key", "value"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedRecord := tt.expectedRecord
			expectedRecord.AddAttrs(tt.expectedAttrs...)

			tt.f(t.Context())

			lastRecord := handler.lastRecord

			// ignore Time and PC diff
			lastRecord.Time = expectedRecord.Time
			lastRecord.PC = expectedRecord.PC

			if !reflect.DeepEqual(expectedRecord, lastRecord) {
				t.Errorf("lastLogHolder = %v, want %v", lastRecord, expectedRecord)
			}
		})
	}
}
