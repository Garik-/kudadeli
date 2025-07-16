package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"kudadeli/bot"
	"kudadeli/config"
	"kudadeli/database"
	"kudadeli/web"

	"golang.org/x/sync/errgroup"
)

// set from ldflags.
var (
	Version = "" //nolint:gochecknoglobals
)

func run(ctx context.Context, cfg *config.Config) error {
	db, err := database.New(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("database.new: %w", err)
	}
	defer db.Close()

	serverHTTP, err := web.New(ctx, cfg.Addr, db)
	if err != nil {
		return fmt.Errorf("failed to create HTTP server: %w", err)
	}

	telebot, err := bot.New(ctx, cfg.Token, db, cfg.AllowUsers)
	if err != nil {
		return fmt.Errorf("telebot new: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		slog.InfoContext(ctx, "starting HTTP server", "address", serverHTTP.Addr)

		return serverHTTP.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()

		return serverHTTP.Shutdown(ctx)
	})

	g.Go(func() error {
		telebot.Start(ctx)

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		telebot.Stop(ctx)

		return nil
	})

	return g.Wait() //nolint:wrapcheck
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.New(Version)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       cfg.LogLevel,
		ReplaceAttr: nil,
		AddSource:   false,
	}))
	slog.SetDefault(logger)

	err := run(ctx, cfg)
	if err != nil {
		logger.ErrorContext(ctx, "run error", "error", err)
	} else {
		logger.InfoContext(ctx, "shutdown")
	}
}
