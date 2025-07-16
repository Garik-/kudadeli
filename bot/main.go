package main

import (
	"context"
	"log/slog"
	"os"

	"kudadeli/config"
	"kudadeli/database"

	_ "modernc.org/sqlite"
)

// set from ldflags.
var (
	Version = "" //nolint:gochecknoglobals
)

func main() {
	ctx := context.Background()
	cfg := config.New(Version)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       cfg.LogLevel,
		ReplaceAttr: nil,
		AddSource:   false,
	}))
	slog.SetDefault(logger)

	db, err := database.New(ctx, cfg.Database)
	if err != nil {
		slog.ErrorContext(ctx, "database.New", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.InfoContext(ctx, "hello world")
}
