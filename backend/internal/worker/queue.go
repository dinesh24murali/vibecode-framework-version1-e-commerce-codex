package worker

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Queue publishes jobs onto the Redis stream consumed by Worker.
type Queue struct {
	rdb *redis.Client
}

// NewQueue returns a Queue backed by rdb.
func NewQueue(rdb *redis.Client) *Queue { return &Queue{rdb: rdb} }

// Enqueue adds a job. id is the idempotency key: a job whose id already completed is skipped.
func (q *Queue) Enqueue(ctx context.Context, id, jobType, payload string) error {
	return q.rdb.XAdd(ctx, &redis.XAddArgs{Stream: StreamName, Values: map[string]any{
		"id": id, "type": jobType, "payload": payload, "attempts": 0,
	}}).Err()
}
