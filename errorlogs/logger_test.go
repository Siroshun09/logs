package errorlogs_test

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"slices"
	"testing"

	"github.com/Siroshun09/logs/errorlogs/v2"
	"github.com/Siroshun09/logs/logmock/v2"
	"github.com/Siroshun09/serrors/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNilLogger(t *testing.T) {
	l := errorlogs.NewNilLogger()
	require.Nil(t, l)

	ctx := t.Context()
	l.Debug(ctx, "test")
	l.Info(ctx, "test")
	l.Warn(ctx, errors.New("test"))
	l.Error(ctx, errors.New("test"))
}

func Test_logger_Debug(t *testing.T) {
	tests := []struct {
		name  string
		opt   *errorlogs.Option
		msg   string
		attrs []slog.Attr
		mock  func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name:  "default option",
			opt:   nil,
			msg:   "test",
			attrs: nil,
			mock: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Debug(ctx, "test")
			},
		},
		{
			name:  "default option with attrs",
			opt:   nil,
			msg:   "test",
			attrs: []slog.Attr{slog.String("key", "value")},
			mock: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Debug(ctx, "test", slog.String("key", "value"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			ctrl := gomock.NewController(t)
			mock := logmock.NewMockLogger(ctrl)
			tt.mock(ctx, mock)

			l := errorlogs.NewLogger(mock, tt.opt)
			l.Debug(ctx, tt.msg, tt.attrs...)
		})
	}
}

func Test_logger_Info(t *testing.T) {
	tests := []struct {
		name  string
		opt   *errorlogs.Option
		msg   string
		attrs []slog.Attr
		mock  func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name:  "default option",
			opt:   nil,
			msg:   "test",
			attrs: nil,
			mock: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Info(ctx, "test")
			},
		},
		{
			name:  "default option with attrs",
			opt:   nil,
			msg:   "test",
			attrs: []slog.Attr{slog.String("key", "value")},
			mock: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Info(ctx, "test", slog.String("key", "value"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			ctrl := gomock.NewController(t)
			mock := logmock.NewMockLogger(ctrl)
			tt.mock(ctx, mock)

			l := errorlogs.NewLogger(mock, tt.opt)
			l.Info(ctx, tt.msg, tt.attrs...)
		})
	}
}

func Test_logger_Warn(t *testing.T) {
	for tt := range testCases() {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			ctrl := gomock.NewController(t)
			mock := logmock.NewMockLogger(ctrl)

			printStackTrace := tt.opt != nil && tt.opt.PrintStackTraceOnWarn
			mock.EXPECT().Warn(ctx, tt.err, createAttrsMatcher(t, tt.err, tt.attrs, printStackTrace, tt.opt))

			l := errorlogs.NewLogger(mock, tt.opt)
			l.Warn(ctx, tt.err, tt.attrs...)
		})
	}
}

func Test_logger_Error(t *testing.T) {
	for tt := range testCases() {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			ctrl := gomock.NewController(t)
			mock := logmock.NewMockLogger(ctrl)
			mock.EXPECT().Error(ctx, tt.err, createAttrsMatcher(t, tt.err, tt.attrs, true, tt.opt))

			l := errorlogs.NewLogger(mock, tt.opt)
			l.Error(ctx, tt.err, tt.attrs...)
		})
	}
}

type testCase struct {
	name  string
	opt   *errorlogs.Option
	err   error
	attrs []slog.Attr
}

func testCases() iter.Seq[testCase] {
	boolMatrix := []bool{true, false}
	stackTraceKeyMatrix := []string{"", errorlogs.DefaultStackTraceAttrKey, "custom_stacktrace_key"}
	errorAttrsFuncs := []func(err error, attrs []slog.Attr) []slog.Attr{
		nil,
		func(err error, attrs []slog.Attr) []slog.Attr {
			return nil // filter all attrs
		},
	}

	opts := []*errorlogs.Option{nil}
	for _, printStackTraceOnWarn := range boolMatrix {
		for _, printCurrentStackTraceIfNotAttached := range boolMatrix {
			for _, stackTraceKey := range stackTraceKeyMatrix {
				for _, errorAttrsFunc := range errorAttrsFuncs {
					opts = append(opts, &errorlogs.Option{
						PrintStackTraceOnWarn:               printStackTraceOnWarn,
						PrintCurrentStackTraceIfNotAttached: printCurrentStackTraceIfNotAttached,
						StackTraceAttrKey:                   stackTraceKey,
						ErrorAttrsFunc:                      errorAttrsFunc,
					})
				}
			}
		}
	}

	errWithAttrs := serrors.New("test_serrors", slog.String("err_key", "err_value"))
	errWithoutAttrs := errors.New("test_errors")
	errs := []error{nil, errWithAttrs, errWithoutAttrs}
	attrsList := [][]slog.Attr{nil, {}, {slog.String("key", "value")}, {slog.String("key_1", "value_1"), slog.String("key_2", "value_2")}}

	return func(yield func(testCase) bool) {
		for _, opt := range opts {
			for _, err := range errs {
				for _, attrs := range attrsList {
					c := testCase{
						name:  fmt.Sprintf("%+v, %+v, %+v", opt, err, attrs),
						opt:   opt,
						err:   err,
						attrs: attrs,
					}
					if !yield(c) {
						return
					}
				}
			}
		}
	}
}

func createAttrsMatcher(t *testing.T, err error, attrs []slog.Attr, printStackTrace bool, opt *errorlogs.Option) gomock.Matcher {
	t.Helper()

	return gomock.Cond(func(got []slog.Attr) bool {
		if len(attrs) != 0 {
			if assert.GreaterOrEqual(t, len(got), len(attrs)) {
				assert.Equal(t, attrs, got[:len(attrs)])
				got = got[len(attrs):]
			}
		}

		expectedStackTraceKey := errorlogs.DefaultStackTraceAttrKey
		if opt != nil && opt.StackTraceAttrKey != "" {
			expectedStackTraceKey = opt.StackTraceAttrKey
		}

		_, hasStackTrace := serrors.GetAttachedStackTrace(err)
		expectStackTrace := false
		if printStackTrace {
			expectStackTrace = hasStackTrace || (opt != nil && opt.PrintCurrentStackTraceIfNotAttached)
		}

		if expectStackTrace {
			if assert.NotEmpty(t, got) {
				assert.Equal(t, expectedStackTraceKey, got[0].Key)
				assert.IsType(t, serrors.StackTrace{}, got[0].Value.Any())
				got = got[1:]
			}
		} else {
			if len(got) != 0 {
				assert.NotEqual(t, expectedStackTraceKey, got[0].Key)
			}
		}

		if opt != nil && opt.ErrorAttrsFunc != nil {
			assert.Empty(t, got)
		} else {
			expectedAttrs := slices.Collect(func(yield func(attr slog.Attr) bool) {
				for _, attr := range serrors.GetAttrs(err) {
					require.True(t, yield(attr))
				}
			})

			if len(expectedAttrs) == 0 {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, expectedAttrs, got)
			}
		}

		return true
	})
}
