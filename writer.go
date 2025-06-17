package logs

import (
	"context"
	"fmt"
	"io"
	"os"
)

func NewStdoutLogger(debug bool) Logger {
	return NewLoggerWithWriter(os.Stdout, true, debug)
}

func NewStderrLogger(debug bool) Logger {
	return NewLoggerWithWriter(os.Stderr, true, debug)
}

func NewLoggerWithWriter(writer io.Writer, printLevel bool, debug bool) Logger {
	return &writerLogger{
		writer:     writer,
		printLevel: printLevel,
		debug:      debug,
	}
}

type writerLogger struct {
	writer     io.Writer
	printLevel bool
	debug      bool
}

func (w *writerLogger) Debug(_ context.Context, msg string) {
	if !w.debug {
		return
	}
	w.println("DEBUG", msg)
}

func (w *writerLogger) Info(_ context.Context, msg string) {
	w.println("INFO", msg)
}

func (w *writerLogger) Warn(_ context.Context, err error) {
	w.println("WARN", err.Error())
}

func (w *writerLogger) Error(_ context.Context, err error) {
	w.println("ERROR", err.Error())
}

func (w *writerLogger) println(level string, msg string) {
	if w.printLevel {
		_, _ = fmt.Fprintln(w.writer, level+":", msg)
	} else {
		_, _ = fmt.Fprintln(w.writer, msg)
	}
}
