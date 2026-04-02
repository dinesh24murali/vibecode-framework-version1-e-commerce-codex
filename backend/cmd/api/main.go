// Command api is the entry point for the e-commerce HTTP API server.
// It loads configuration from environment variables, builds the HTTP server
// using the Builder pattern, and handles graceful shutdown on SIGINT/SIGTERM.
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/cache"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/config"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/db"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/server"
)

func main() {
	// Structured JSON logging — all log output uses log/slog so that
	// CloudWatch can parse log lines as JSON.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load and validate configuration. The Adapter calls log.Fatal with a
	// clear message if any required env var is absent.
	cfg := config.NewEnvAdapter().Load()

	ctx := context.Background()

	// Open the PostgreSQL connection pool. Exit immediately with a clear error
	// if the database is unreachable — fast-fail is better than a half-started server.
	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("startup: database unreachable: %v", err)
	}
	defer pool.Close()
	slog.Info("database connected")

	// Connect to Redis. Same fast-fail behaviour.
	redisClient, err := cache.NewClient(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("startup: redis unreachable: %v", err)
	}
	defer redisClient.Close()
	slog.Info("redis connected")

	// Build the HTTP server using the Builder pattern. Defaults are
	// production-safe timeouts; override here if the environment requires it.
	srv := server.NewServerBuilder(cfg).
		WithPool(pool).
		WithRedis(redisClient).
		Build()

	// Start the server in a goroutine so the main goroutine can listen for
	// OS signals for graceful shutdown.
	idleConnsClosed := make(chan struct{})
	go func() {
		slog.Info("server starting", "addr", cfg.Addr())

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Block until SIGINT or SIGTERM is received.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("shutdown signal received", "signal", sig.String())

	// Allow up to ShutdownTimeout for in-flight requests to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), server.ShutdownTimeout())
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	close(idleConnsClosed)
	<-idleConnsClosed

	slog.Info("server stopped")
}
