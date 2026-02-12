package plain

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/Siroshun09/logs/v2"
)

type Option struct {
	Debug      bool
	PrintLevel bool
	PrintAttrs bool
}

func NewPlainLogger(writer io.Writer, option *Option) logs.Logger {
	opt := Option{}
	if option != nil {
		opt = *option
	}

	return &writerLogger{
		writer: writer,
		opt:    opt,
	}
}

type writerLogger struct {
	writer io.Writer
	opt    Option
}

func (w *writerLogger) Debug(_ context.Context, msg string, attrs ...slog.Attr) {
	if !w.opt.Debug {
		return
	}
	w.println("DEBUG", msg, attrs...)
}

func (w *writerLogger) Info(_ context.Context, msg string, attrs ...slog.Attr) {
	w.println("INFO", msg, attrs...)
}

func (w *writerLogger) Warn(_ context.Context, err error, attrs ...slog.Attr) {
	w.println("WARN", err.Error(), attrs...)
}

func (w *writerLogger) Error(_ context.Context, err error, attrs ...slog.Attr) {
	w.println("ERROR", err.Error(), attrs...)
}

func (w *writerLogger) println(level string, msg string, attrs ...slog.Attr) {
	if w.opt.PrintLevel {
		_, _ = fmt.Fprint(w.writer, level+": ")
	}

	_, _ = fmt.Fprint(w.writer, msg)

	if w.opt.PrintAttrs {
		for _, attr := range attrs {
			_, _ = fmt.Fprintf(w.writer, " %s=%v", attr.Key, attr.Value)
		}
	}

	_, _ = fmt.Fprintln(w.writer)
}
