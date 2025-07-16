package main

import (
	"log/slog"
	"os"
)

// set from ldflags.
var (
	Version = "" //nolint:gochecknoglobals
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("hello world")
}
