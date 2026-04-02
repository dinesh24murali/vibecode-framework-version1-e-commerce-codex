package cache_test

import (
	"context"
	"os"
	"testing"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/cache"
)

func TestNewClient_Integration(t *testing.T) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL not set; skipping integration test")
	}

	ctx := context.Background()
	client, err := cache.NewClient(ctx, url)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}
