# logs

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/Siroshun09/logs)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Siroshun09/logs/ci.yml?branch=main)
![GitHub](https://img.shields.io/github/license/Siroshun09/logs)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Siroshun09/logs)

A Go library that provides a simple interface and functions for logging.

## Requirements

- Go 1.24+

## Installation

```shell
go get github.com/Siroshun09/logs/v2
```

The following sub-modules are also available:

```shell
# A plain-text logger that writes to an io.Writer is provided by the v2 module above (package plain).

# A logger that appends stack traces / attributes from errors created with Siroshun09/serrors.
go get github.com/Siroshun09/logs/errorlogs/v2

# A generated GoMock implementation of the Logger interface, for tests.
go get github.com/Siroshun09/logs/logmock/v2
```

## Usage

This library provides a common `Logger` interface and simple context-based utility functions.
By default, it works with the standard `log/slog` package, and you can swap the output destination or the default
logger to suit your needs.

The `Logger` interface is:

```go
package logs

import (
	"context"
	"log/slog"
)

type Logger interface {
    Debug(ctx context.Context, msg string, attrs ...slog.Attr)
    Info(ctx context.Context, msg string, attrs ...slog.Attr)
    Warn(ctx context.Context, err error, attrs ...slog.Attr)
    Error(ctx context.Context, err error, attrs ...slog.Attr)
}
```

`Debug`/`Info` take a message, while `Warn`/`Error` take an `error`. Every method also accepts optional
`slog.Attr` values to attach structured attributes to the log entry.

### Quick start

```go
package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Siroshun09/logs/v2"
)

func main() {
	ctx := context.Background()

	// Log with the default logger (internally uses slog.Default)
	logs.Info(ctx, "hello")
	logs.Debug(ctx, "debug message") // whether this appears depends on slog's level settings

	// Attach structured attributes
	logs.Info(ctx, "user logged in", slog.String("user_id", "123"))

	// Log errors at warn/error levels
	logs.Warn(ctx, errors.New("something not critical"))
	logs.Error(ctx, errors.New("something went wrong"), slog.Int("code", 500))
}
```

### Get and set the default logger

```go
package main

import (
	"github.com/Siroshun09/logs/v2"
)

func main() {
	// Current default logger (when SetDefault is not called, it's based on slog.Default)
	logger := logs.Default()

	// Set any Logger as the default (re-setting the same here as an example)
	logs.SetDefault(logger)
}
```

### Integrate with slog

You can pass a `slog.Logger` and use it as this library's `Logger`.

```go
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/Siroshun09/logs/v2"
)

func main() {
	ctx := context.Background()

	// Prepare slog with a JSON handler
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	s := slog.New(h)

	// Convert to this library's Logger and set as default
	logs.SetDefault(logs.NewLoggerWithSlog(s))

	logs.Info(ctx, "info message")
	logs.Warn(ctx, errors.New("warn!"))
}
```

If you pass nil to `logs.NewLoggerWithSlog(nil)`, it uses `slog.Default()`.

### Simple plain-text logger

The `plain` package provides a simple `Logger` that writes plain text to any `io.Writer`.

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Siroshun09/logs/v2"
	"github.com/Siroshun09/logs/v2/plain"
)

func main() {
	ctx := context.Background()

	l := plain.NewPlainLogger(os.Stdout, &plain.Option{
		Debug:      true, // print Debug logs (otherwise they are skipped)
		PrintLevel: true, // prefix each line with the level, e.g. "INFO: ..."
		PrintAttrs: true, // append attributes as " key=value"
	})
	logs.SetDefault(l)

	logs.Debug(ctx, "debug visible")                      // printed because Debug=true
	logs.Info(ctx, "info message", slog.String("k", "v")) // prints: "INFO: info message k=v"
}
```

`NewPlainLogger` accepts a nil `*Option`, in which case all options default to `false`
(Debug logs skipped, no level prefix, no attributes).

### Append stack traces and attributes from errors

The `errorlogs` package wraps another `Logger` and, on `Warn`/`Error`, appends the stack trace and attributes
attached to errors created with [Siroshun09/serrors](https://github.com/Siroshun09/serrors).

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Siroshun09/logs/v2"
	"github.com/Siroshun09/logs/errorlogs/v2"
	"github.com/Siroshun09/serrors/v2"
)

func main() {
	ctx := context.Background()

	// Wrap any Logger (here a slog-backed JSON logger).
	base := logs.NewLoggerWithSlog(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	l := errorlogs.NewLogger(base, &errorlogs.Option{
		// Print the stack trace on Warn too (it is always printed on Error).
		PrintStackTraceOnWarn: true,
		// If the error has no attached stack trace, capture the current one.
		PrintCurrentStackTraceIfNotAttached: true,
		// The attribute key for the stack trace (defaults to errorlogs.DefaultStackTraceAttrKey = "stack_trace").
		StackTraceAttrKey: "stack_trace",
		// Optionally transform the attributes extracted from the error.
		ErrorAttrsFunc: func(err error, attrs []slog.Attr) []slog.Attr {
			return attrs
		},
	})

	// The stack trace and attributes attached to the error are logged automatically.
	l.Error(ctx, serrors.New("something went wrong", slog.String("op", "doWork")))
}
```

`NewLogger` accepts a nil `*Option`, in which case stack traces are printed only on `Error`, no current
stack trace is captured when the error has none, and the default attribute key is used.

### Switch logger via context

If you want to temporarily swap the logger per request or operation, use `WithContext`/`FromContext` or the
package-level functions.

```go
package main

import (
	"context"
	"os"

	"github.com/Siroshun09/logs/v2"
	"github.com/Siroshun09/logs/v2/plain"
)

func handler(ctx context.Context) {
	// Use a different logger only for this logic
	l := plain.NewPlainLogger(os.Stdout, &plain.Option{Debug: true, PrintLevel: true})
	ctx = logs.WithContext(ctx, l)

	// Package-level functions internally call FromContext(ctx)
	logs.Info(ctx, "in handler")
}

func main() {
	handler(context.Background())
}
```

`FromContext(ctx)` returns `logs.Default()` when no logger is attached to the context.

### Mocking in tests

The `logmock` module provides a generated GoMock implementation of `Logger`.

```go
ctrl := gomock.NewController(t)
mock := logmock.NewMockLogger(ctrl)
mock.EXPECT().Info(ctx, "hello", slog.String("key", "value"))

mock.Info(ctx, "hello", slog.String("key", "value"))
```

## License

This project is under the Apache License version 2.0. Please see [LICENSE](LICENSE) for more info.

Copyright © 2024-2026, Siroshun09
