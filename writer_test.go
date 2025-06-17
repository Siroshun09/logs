package logs

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewStdoutLogger(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
		want  Logger
	}{
		{
			name:  "debug true",
			debug: true,
			want: &writerLogger{
				writer:     os.Stdout,
				printLevel: true,
				debug:      true,
			},
		},
		{
			name:  "debug false",
			debug: false,
			want: &writerLogger{
				writer:     os.Stdout,
				printLevel: true,
				debug:      false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStdoutLogger(tt.debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStdoutLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStderrLogger(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
		want  Logger
	}{
		{
			name:  "debug true",
			debug: true,
			want: &writerLogger{
				writer:     os.Stderr,
				printLevel: true,
				debug:      true,
			},
		},
		{
			name:  "debug false",
			debug: false,
			want: &writerLogger{
				writer:     os.Stderr,
				printLevel: true,
				debug:      false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStderrLogger(tt.debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStderrLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLoggerWithWriter(t *testing.T) {
	tests := []struct {
		name       string
		writer     io.Writer
		printLevel bool
		debug      bool
		want       Logger
	}{
		{
			name:       "printLevel: true | debug: true",
			writer:     io.Discard,
			printLevel: true,
			debug:      true,
			want: &writerLogger{
				writer:     io.Discard,
				printLevel: true,
				debug:      true,
			},
		},
		{
			name:       "printLevel: false | debug: false",
			writer:     io.Discard,
			printLevel: false,
			debug:      false,
			want: &writerLogger{
				writer:     io.Discard,
				printLevel: false,
				debug:      false,
			},
		},
		{
			name:       "printLevel: true | debug: false",
			writer:     io.Discard,
			printLevel: true,
			debug:      false,
			want: &writerLogger{
				writer:     io.Discard,
				printLevel: true,
				debug:      false,
			},
		},
		{
			name:       "printLevel: false | debug: true",
			writer:     io.Discard,
			printLevel: false,
			debug:      true,
			want: &writerLogger{
				writer:     io.Discard,
				printLevel: false,
				debug:      true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLoggerWithWriter(tt.writer, tt.printLevel, tt.debug)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLoggerWithWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writerLogger_Debug(t *testing.T) {
	tests := []struct {
		name         string
		printLevel   bool
		debug        bool
		msg          string
		expectedLogs []string
	}{
		{
			name:         "printLevel: true | debug: true",
			printLevel:   true,
			debug:        true,
			msg:          "debug log",
			expectedLogs: []string{"DEBUG: debug log"},
		},
		{
			name:         "printLevel: false | debug: false",
			printLevel:   false,
			debug:        false,
			msg:          "debug log",
			expectedLogs: []string{},
		},
		{
			name:         "printLevel: true | debug: false",
			printLevel:   true,
			debug:        false,
			msg:          "debug log",
			expectedLogs: []string{},
		},
		{
			name:         "printLevel: false | debug: true",
			printLevel:   false,
			debug:        true,
			msg:          "debug log",
			expectedLogs: []string{"debug log"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := strings.Builder{}
			w := &writerLogger{
				writer:     &buf,
				printLevel: tt.printLevel,
				debug:      tt.debug,
			}
			w.Debug(context.Background(), tt.msg)
			got := strings.Split(buf.String(), "\n")
			if 0 < len(got) {
				got = got[:len(got)-1] // remove trailing line
			}
			if !reflect.DeepEqual(got, tt.expectedLogs) {
				t.Errorf("Debug() = %v, want %v", got, tt.expectedLogs)
			}
		})
	}
}

func Test_writerLogger_Info(t *testing.T) {
	tests := []struct {
		name         string
		printLevel   bool
		msg          string
		expectedLogs []string
	}{
		{
			name:         "printLevel: true",
			printLevel:   true,
			msg:          "info log",
			expectedLogs: []string{"INFO: info log"},
		},
		{
			name:         "printLevel: false",
			printLevel:   false,
			msg:          "info log",
			expectedLogs: []string{"info log"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := strings.Builder{}
			w := &writerLogger{
				writer:     &buf,
				printLevel: tt.printLevel,
			}
			w.Info(context.Background(), tt.msg)
			got := strings.Split(buf.String(), "\n")
			if 0 < len(got) {
				got = got[:len(got)-1] // remove trailing line
			}
			if !reflect.DeepEqual(got, tt.expectedLogs) {
				t.Errorf("Info() = %v, want %v", got, tt.expectedLogs)
			}
		})
	}
}

func Test_writerLogger_Warn(t *testing.T) {
	tests := []struct {
		name         string
		printLevel   bool
		err          error
		expectedLogs []string
	}{
		{
			name:         "printLevel: true",
			printLevel:   true,
			err:          fmt.Errorf("warning log"),
			expectedLogs: []string{"WARN: warning log"},
		},
		{
			name:         "printLevel: false",
			printLevel:   false,
			err:          fmt.Errorf("warning log"),
			expectedLogs: []string{"warning log"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := strings.Builder{}
			w := &writerLogger{
				writer:     &buf,
				printLevel: tt.printLevel,
			}
			w.Warn(context.Background(), tt.err)
			got := strings.Split(buf.String(), "\n")
			if 0 < len(got) {
				got = got[:len(got)-1] // remove trailing line
			}
			if !reflect.DeepEqual(got, tt.expectedLogs) {
				t.Errorf("Warn() = %v, want %v", got, tt.expectedLogs)
			}
		})
	}
}

func Test_writerLogger_Error(t *testing.T) {
	tests := []struct {
		name         string
		printLevel   bool
		err          error
		expectedLogs []string
	}{
		{
			name:         "printLevel: true",
			printLevel:   true,
			err:          fmt.Errorf("error log"),
			expectedLogs: []string{"ERROR: error log"},
		},
		{
			name:         "printLevel: false",
			printLevel:   false,
			err:          fmt.Errorf("error log"),
			expectedLogs: []string{"error log"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := strings.Builder{}
			w := &writerLogger{
				writer:     &buf,
				printLevel: tt.printLevel,
			}
			w.Error(context.Background(), tt.err)
			got := strings.Split(buf.String(), "\n")
			if 0 < len(got) {
				got = got[:len(got)-1] // remove trailing line
			}
			if !reflect.DeepEqual(got, tt.expectedLogs) {
				t.Errorf("Error() = %v, want %v", got, tt.expectedLogs)
			}
		})
	}
}
