package services_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// Integration: claim, in-flight conflict, replay of the stored response, and abort.
func TestRedisIdempotencyStore_Integration(t *testing.T) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL not set")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	rdb := redis.NewClient(opts)
	defer func() { _ = rdb.Close() }()
	ctx := context.Background()
	store := services.NewRedisIdempotencyStore(rdb)
	key := "idem:test:" + uuid.NewString()
	defer rdb.Del(ctx, key)

	if _, done, e := store.Begin(ctx, key); e != nil || done {
		t.Fatalf("first Begin: done=%v err=%v", done, e)
	}
	if _, _, e := store.Begin(ctx, key); !errors.Is(e, services.ErrIdempotencyInFlight) {
		t.Fatalf("concurrent Begin: %v", e)
	}
	if err = store.Complete(ctx, key, []byte(`{"a":1}`)); err != nil {
		t.Fatal(err)
	}
	v, done, err := store.Begin(ctx, key)
	if err != nil || !done || string(v) != `{"a":1}` {
		t.Fatalf("replay: %s done=%v err=%v", v, done, err)
	}

	key2 := key + ":abort"
	defer rdb.Del(ctx, key2)
	if _, _, err := store.Begin(ctx, key2); err != nil {
		t.Fatal(err)
	}
	store.Abort(ctx, key2)
	if _, done, err := store.Begin(ctx, key2); err != nil || done {
		t.Fatalf("after abort: done=%v err=%v", done, err)
	}
}
