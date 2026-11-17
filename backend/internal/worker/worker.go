package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Stream and group names shared by all worker replicas.
const (
	StreamName     = "jobs"
	GroupName      = "workers"
	DeadLetterName = "jobs:dead-letter"
	MaxRetries     = 3
	doneKeyPrefix  = "job:done:"
	doneKeyTTL     = 7 * 24 * time.Hour
)

// Backoff returns the retry delay for the given attempt (1s, doubling, cap 60s).
func Backoff(attempt int) time.Duration {
	d := time.Second
	for i := 1; i < attempt; i++ {
		d *= 2
		if d >= 60*time.Second {
			return 60 * time.Second
		}
	}
	return d
}

// Worker consumes jobs from a Redis Stream and dispatches them.
type Worker struct {
	rdb        *redis.Client
	dispatcher *Dispatcher
	backoff    func(int) time.Duration
	consumer   string
	block      time.Duration
	minIdle    time.Duration
}

// New creates a Worker. consumer must be unique per replica.
func New(rdb *redis.Client, d *Dispatcher, consumer string) *Worker {
	return &Worker{
		rdb: rdb, dispatcher: d, consumer: consumer,
		block: 2 * time.Second, minIdle: time.Minute, backoff: Backoff,
	}
}

// Run blocks, processing jobs until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	err := w.rdb.XGroupCreateMkStream(ctx, StreamName, GroupName, "0").Err()
	if err != nil && !strings.HasPrefix(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("worker: create group: %w", err)
	}
	for ctx.Err() == nil {
		w.reclaim(ctx)
		streams, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group: GroupName, Consumer: w.consumer,
			Streams: []string{StreamName, ">"}, Count: 10, Block: w.block,
		}).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) || ctx.Err() != nil {
				continue
			}
			slog.Error("worker: read", "err", err)
			time.Sleep(time.Second)
			continue
		}
		for _, s := range streams {
			for _, m := range s.Messages {
				w.process(ctx, m)
			}
		}
	}
	return nil
}

// reclaim takes over messages left pending by crashed consumers.
func (w *Worker) reclaim(ctx context.Context) {
	msgs, _, err := w.rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream: StreamName, Group: GroupName, Consumer: w.consumer,
		MinIdle: w.minIdle, Start: "0", Count: 10,
	}).Result()
	if err != nil {
		return
	}
	for _, m := range msgs {
		w.process(ctx, m)
	}
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

func (w *Worker) process(ctx context.Context, m redis.XMessage) {
	attempts, _ := strconv.Atoi(str(m.Values["attempts"]))
	job := Job{
		ID: str(m.Values["id"]), Type: str(m.Values["type"]),
		Payload: str(m.Values["payload"]), Attempts: attempts,
	}
	if job.ID == "" {
		job.ID = m.ID
	}

	doneKey := doneKeyPrefix + job.ID
	if n, _ := w.rdb.Exists(ctx, doneKey).Result(); n > 0 {
		w.ack(ctx, m.ID) // already completed: replay is a no-op
		return
	}

	err := w.dispatcher.Dispatch(ctx, job)
	if err == nil {
		w.rdb.Set(ctx, doneKey, "1", doneKeyTTL)
		w.ack(ctx, m.ID)
		return
	}
	slog.Error("worker: job failed", "id", job.ID, "type", job.Type, "attempt", job.Attempts, "err", err)

	if job.Attempts >= MaxRetries {
		w.rdb.XAdd(ctx, &redis.XAddArgs{Stream: DeadLetterName, Values: map[string]any{
			"id": job.ID, "type": job.Type, "payload": job.Payload,
			"attempts": job.Attempts, "error": err.Error(),
		}})
		w.ack(ctx, m.ID)
		return
	}
	// Re-add after backoff, then ack the original. If we crash first the
	// original stays pending and is reclaimed via XAUTOCLAIM.
	next := job.Attempts + 1
	origID := m.ID
	time.AfterFunc(w.backoff(next), func() {
		bg := context.Background()
		if err := w.rdb.XAdd(bg, &redis.XAddArgs{Stream: StreamName, Values: map[string]any{
			"id": job.ID, "type": job.Type, "payload": job.Payload, "attempts": next,
		}}).Err(); err != nil {
			slog.Error("worker: requeue", "id", job.ID, "err", err)
			return
		}
		w.ack(bg, origID)
	})
}

func (w *Worker) ack(ctx context.Context, id string) {
	if err := w.rdb.XAck(ctx, StreamName, GroupName, id).Err(); err != nil {
		slog.Error("worker: ack", "id", id, "err", err)
	}
}
