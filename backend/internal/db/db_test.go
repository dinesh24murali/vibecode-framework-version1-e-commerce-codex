package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/db"
)

func TestOpen_Integration(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping integration test")
	}

	ctx := context.Background()
	pool, err := db.Open(ctx, url)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Ping after Open: %v", err)
	}

	stat := pool.Stat()
	if stat.MaxConns() != 20 {
		t.Errorf("MaxConns = %d, want 20", stat.MaxConns())
	}
}
