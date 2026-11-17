package worker

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestBackoff(t *testing.T) {
	want := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second}
	for i, w := range want {
		if got := Backoff(i + 1); got != w {
			t.Errorf("Backoff(%d)=%v want %v", i+1, got, w)
		}
	}
	if Backoff(20) != 60*time.Second {
		t.Error("backoff must cap at 60s")
	}
}

func TestDispatcher(t *testing.T) {
	d := NewDispatcher()
	var n int32
	d.Register("a", HandlerFunc(func(context.Context, Job) error { atomic.AddInt32(&n, 1); return nil }))
	if err := d.Dispatch(context.Background(), Job{Type: "a"}); err != nil || n != 1 {
		t.Fatalf("dispatch: %v n=%d", err, n)
	}
	if err := d.Dispatch(context.Background(), Job{Type: "missing"}); err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func testRedis(t *testing.T) *redis.Client {
	t.Helper()
	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL not set")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	rdb := redis.NewClient(opts)
	ctx := context.Background()
	rdb.Del(ctx, StreamName, DeadLetterName, "job:done:ok1", "job:done:bad1")
	return rdb
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	for i := 0; i < 100; i++ {
		if cond() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("condition not met in time")
}

func TestWorker_SuccessAndDeadLetter(t *testing.T) {
	rdb := testRedis(t)
	defer rdb.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var okCalls, badCalls int32
	d := NewDispatcher()
	d.Register("ok", HandlerFunc(func(context.Context, Job) error { atomic.AddInt32(&okCalls, 1); return nil }))
	d.Register("bad", HandlerFunc(func(context.Context, Job) error {
		atomic.AddInt32(&badCalls, 1)
		return errors.New("boom")
	}))

	w := New(rdb, d, "test")
	w.block = 200 * time.Millisecond
	w.backoff = func(int) time.Duration { return 10 * time.Millisecond }
	go func() { _ = w.Run(ctx) }()

	// Same job id twice: idempotency means the handler runs once.
	for i := 0; i < 2; i++ {
		rdb.XAdd(ctx, &redis.XAddArgs{Stream: StreamName, Values: map[string]any{"id": "ok1", "type": "ok"}})
	}
	rdb.XAdd(ctx, &redis.XAddArgs{Stream: StreamName, Values: map[string]any{"id": "bad1", "type": "bad"}})

	waitFor(t, func() bool { n, _ := rdb.XLen(ctx, DeadLetterName).Result(); return n == 1 })
	if got := atomic.LoadInt32(&okCalls); got != 1 {
		t.Errorf("ok handler ran %d times, want 1", got)
	}
	if got := atomic.LoadInt32(&badCalls); got != MaxRetries+1 {
		t.Errorf("bad handler ran %d times, want %d", got, MaxRetries+1)
	}
}
