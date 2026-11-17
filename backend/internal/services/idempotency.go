package services

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrIdempotencyInFlight is returned while the first request with the same key is still running.
var ErrIdempotencyInFlight = errors.New("request with this idempotency key is in progress")

// IdempotencyStore remembers the response of the first request for a key.
type IdempotencyStore interface {
	// Begin claims key. When a response is already stored it is returned with done=true.
	// ErrIdempotencyInFlight means another request holds the claim.
	Begin(ctx context.Context, key string) (cached []byte, done bool, err error)
	Complete(ctx context.Context, key string, response []byte) error
	// Abort releases the claim after a failure so the client can retry.
	Abort(ctx context.Context, key string)
}

const (
	idemPending = "pending"
	idemLockTTL = 30 * time.Second
	idemTTL     = 24 * time.Hour
)

type redisIdempotency struct{ rdb *redis.Client }

// NewRedisIdempotencyStore returns an IdempotencyStore backed by Redis.
func NewRedisIdempotencyStore(rdb *redis.Client) IdempotencyStore { return &redisIdempotency{rdb: rdb} }

func (s *redisIdempotency) Begin(ctx context.Context, key string) ([]byte, bool, error) {
	ok, err := s.rdb.SetArgs(ctx, key, idemPending, redis.SetArgs{Mode: "NX", TTL: idemLockTTL}).Result()
	if errors.Is(err, redis.Nil) { // key already exists
		ok, err = "", nil
	}
	if err != nil {
		return nil, false, err
	}
	if ok == "OK" {
		return nil, false, nil
	}
	v, err := s.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, ErrIdempotencyInFlight // expired between calls: let the client retry
	}
	if err != nil {
		return nil, false, err
	}
	if string(v) == idemPending {
		return nil, false, ErrIdempotencyInFlight
	}
	return v, true, nil
}

func (s *redisIdempotency) Complete(ctx context.Context, key string, response []byte) error {
	return s.rdb.Set(ctx, key, response, idemTTL).Err()
}

func (s *redisIdempotency) Abort(ctx context.Context, key string) {
	s.rdb.Del(ctx, key)
}
