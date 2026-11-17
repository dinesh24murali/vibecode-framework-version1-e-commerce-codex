// Command worker runs the background job consumer (Redis Streams).
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/cache"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/db"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/worker"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/worker/handlers"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	dbURL, redisURL := os.Getenv("DATABASE_URL"), os.Getenv("REDIS_URL")
	if dbURL == "" || redisURL == "" {
		log.Fatal("startup: DATABASE_URL and REDIS_URL are required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.Open(ctx, dbURL)
	if err != nil {
		log.Fatalf("startup: database unreachable: %v", err)
	}
	defer pool.Close()

	rdb, err := cache.NewClient(ctx, redisURL)
	if err != nil {
		log.Fatalf("startup: redis unreachable: %v", err)
	}
	defer rdb.Close()

	d := worker.NewDispatcher()
	handlers.RegisterAll(d)

	host, _ := os.Hostname()
	slog.Info("worker started", "consumer", host)
	if err := worker.New(rdb, d, host).Run(ctx); err != nil {
		log.Fatalf("worker: %v", err)
	}
	slog.Info("worker stopped")
}
