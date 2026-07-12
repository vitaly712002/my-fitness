// Command api is the backend's HTTP server entrypoint.
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

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vitaly712002/my-fitness/backend/internal/config"
	"github.com/vitaly712002/my-fitness/backend/internal/handler"
	"github.com/vitaly712002/my-fitness/backend/internal/repository"
)

func main() {
	cfg := config.Load()

	// JSON handler: structured logs are what a real deployment (and slog
	// itself) is built around — easy to grep/parse, unlike printf-style logs.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.ParseLogLevel(),
	}))
	slog.SetDefault(logger)

	// A short-lived context just for startup (connecting to Postgres and
	// verifying it), separate from the server's own long-running lifetime.
	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelStartup()

	pool, err := pgxpool.New(startupCtx, cfg.DatabaseURL())
	if err != nil {
		logger.Error("failed to create postgres pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Fail fast at startup rather than accepting traffic against a database
	// we can't actually reach.
	if err := pool.Ping(startupCtx); err != nil {
		logger.Error("failed to ping postgres", "error", err)
		os.Exit(1)
	}

	queries := repository.New(pool)
	router := handler.NewRouter(queries, logger)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run the server in its own goroutine so the main goroutine is free to
	// wait for an OS shutdown signal below.
	go func() {
		logger.Info("server listening", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Block until Ctrl+C (SIGINT) or a process manager's SIGTERM, then shut
	// down cleanly instead of dropping in-flight requests.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Info("shutting down")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
