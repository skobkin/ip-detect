// Package main boots the ip-detect HTTP server.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"git.skobk.in/skobkin/ip-detect/internal/config"
	"git.skobk.in/skobkin/ip-detect/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)

		os.Exit(1)
	}

	logger := newLogger(cfg.Logging)

	srv, err := server.New(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize server", "error", err)

		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err = srv.Run(ctx)

	stop()

	if err != nil {
		logger.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}

func newLogger(cfg config.LoggingConfig) *slog.Logger {
	handlerOpts := &slog.HandlerOptions{Level: cfg.Level}

	switch cfg.Format {
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))
	}
}
