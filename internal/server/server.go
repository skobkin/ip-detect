// Package server wires HTTP handlers and runs the HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"git.skobk.in/skobkin/ip-detect/internal/config"
)

// App wraps the HTTP server lifecycle.
type App struct {
	cfg        config.Config
	logger     *slog.Logger
	httpServer *http.Server
}

// New constructs a server with routes configured.
func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	handler, err := newHandler(cfg, logger)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           handler,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		MaxHeaderBytes:    cfg.Server.MaxHeaderBytes,
	}

	return &App{cfg: cfg, logger: logger, httpServer: srv}, nil
}

// Run starts the HTTP server and blocks until shutdown.
func (a *App) Run(ctx context.Context) error {
	serverErr := make(chan error, 1)

	go func() {
		a.logger.Info("server listening", "addr", a.cfg.Server.Addr)

		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err

			return
		}

		serverErr <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(ctx, a.cfg.Server.ShutdownTimeout)
		defer cancel()

		a.logger.Info("shutting down")

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}

		return <-serverErr
	case err := <-serverErr:
		return err
	}
}
