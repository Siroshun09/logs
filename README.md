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
go get github.com/Siroshun09/logs
```

## Usage

This library provides a common Logger interface and simple context-based utility functions.
By default, it works with the standard log/slog package, and you can swap the output destination or the default logger
to suit your needs.

### Quick start

```
package main

import (
    "context"
    "errors"
    "github.com/Siroshun09/logs"
)

func main() {
    ctx := context.Background()

    // Log with the default logger (internally uses slog.Default)
    logs.Info(ctx, "hello")
    logs.Debug(ctx, "debug message") // whether this appears depends on slog's level settings

    // Log errors at warn/error levels
    logs.Warn(ctx, errors.New("something not critical"))
    logs.Error(ctx, errors.New("something went wrong"))
}
```

### Get and set the default logger

```
package main

import (
    "github.com/Siroshun09/logs"
)

func main() {
    // Current default logger (when SetDefault is not called, it's based on slog.Default)
    logger := logs.Default()

    // Set any Logger as the default (re-setting the same here as an example)
    logs.SetDefault(logger)
}
```

### Integrate with slog

You can pass a slog.Logger and use it as this library's Logger.

```
package main

import (
    "context"
    "errors"
    "log/slog"
    "os"
    "github.com/Siroshun09/logs"
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

### Simple writer-based logger

You can also use a simple plain-text Logger that writes to standard output/error.

```
package main

import (
    "context"
    "os"
    "github.com/Siroshun09/logs"
)

func main() {
    ctx := context.Background()

    // Write INFO and above to standard output.
    // Arguments: NewLoggerWithWriter(writer, printLevel, debug)
    l := logs.NewLoggerWithWriter(os.Stdout, true /*printLevel*/, true /*debug*/)
    logs.SetDefault(l)

    logs.Debug(ctx, "debug visible") // printed because debug=true
    logs.Info(ctx, "info message")   // prints like: "INFO: info message"
}
```

Helpers are also provided:

- `logs.NewStdoutLogger(debug bool)`
- `logs.NewStderrLogger(debug bool)`

### Switch logger via context

If you want to temporarily swap the logger per request or operation, use `WithContext`/`FromContext` or the
package-level functions.

```
package main

import (
    "context"
    "github.com/Siroshun09/logs"
)

func handler(ctx context.Context) {
    // Use a different logger only for this logic
    l := logs.NewStdoutLogger(true)
    ctx = logs.WithContext(ctx, l)

    // Package-level functions internally call FromContext(ctx)
    logs.Info(ctx, "in handler")
}

func main() {
    handler(context.Background())
}
```

`FromContext(ctx)` returns `logs.Default()` when no logger is attached to the context.

## License

This project is under the Apache License version 2.0. Please see [LICENSE](LICENSE) for more info.

Copyright © 2024-2025, Siroshun09
