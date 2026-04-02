// Command migrate is a one-shot CLI for running goose database migrations.
// It reads DATABASE_URL from the environment and accepts a single argument:
// "up", "down", or "status".
//
// Usage:
//
//	make migrate-up
//	make migrate-down
//	make migrate-status
package main

import (
	"context"
	"log"
	"os"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/db"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("migrate: usage: migrate <up|down|status>")
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("migrate: DATABASE_URL is required")
	}

	ctx := context.Background()
	var err error

	switch os.Args[1] {
	case "up":
		err = db.RunUp(ctx, url)
	case "down":
		err = db.RunDown(ctx, url)
	case "status":
		err = db.PrintStatus(ctx, url)
	default:
		log.Fatalf("migrate: unknown command %q (use up, down, or status)", os.Args[1])
	}

	if err != nil {
		log.Fatalf("migrate: %v", err)
	}
}
