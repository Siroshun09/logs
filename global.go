package logs

import "sync/atomic"

var globalLogger atomic.Pointer[Logger]

func init() {
	logger := NewLoggerWithSlog(nil)
	globalLogger.Store(&logger)
}

// Default returns a Logger set by SetDefault, or default Logger instance that uses slog.Default.
func Default() Logger {
	return *globalLogger.Load()
}

// SetDefault makes Default return the specified Logger.
//
// If the given Logger is nil, this function triggers a panic.
func SetDefault(logger Logger) {
	if logger == nil {
		panic("default logger cannot be nil")
	}
	globalLogger.Store(&logger)
}
