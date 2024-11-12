package logs

import (
	"log/slog"
	"testing"
)

func TestSetDefault(t *testing.T) {
	logger := NewLoggerWithSlog(&slog.Logger{})
	SetDefault(logger)
	if Default() != logger {
		t.Errorf("SetDefault failed")
	}

}

func TestSetDefaultPanic(t *testing.T) {
	defer func() {
		err := recover()
		if err != "default logger cannot be nil" {
			t.Errorf("got %v\nwant %v", err, "default logger cannot be nil")
		}
	}()

	SetDefault(nil)
}
