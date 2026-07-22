// Command server runs the CleanLeaf website API.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cncmarella/cleanleaf/api/internal/config"
	"github.com/cncmarella/cleanleaf/api/internal/handler"
	"github.com/cncmarella/cleanleaf/api/internal/mailer"
)

// version is stamped at build time: -ldflags "-X main.version=$(git rev-parse --short HEAD)".
var version = "dev"

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if err := run(log); err != nil {
		log.Error("server exited", "error", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	var m mailer.Mailer
	if cfg.ResendAPIKey != "" {
		m = mailer.NewResend(cfg.ResendAPIKey)
	} else {
		log.Warn("RESEND_API_KEY not set; contact enquiries will only be logged")
		m = mailer.NewConsole(log)
	}

	srv, routes := handler.New(cfg, log, m, version)

	// Signal handling is installed before the listener so a Ctrl-C during
	// startup still unwinds cleanly.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go reapLimiter(ctx, srv)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", httpServer.Addr, "env", cfg.Env, "version", version)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Info("shutdown complete")
	return nil
}

// reapLimiter drops expired rate-limit buckets so memory stays flat.
func reapLimiter(ctx context.Context, srv *handler.Server) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			srv.Limiter().Reap()
		}
	}
}
