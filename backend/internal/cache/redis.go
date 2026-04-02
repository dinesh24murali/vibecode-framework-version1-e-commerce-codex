// Package cache manages the Redis client lifecycle.
package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewClient parses url, constructs a go-redis/v9 client, and pings Redis
// to verify connectivity. Returns an error if the URL is invalid or Redis
// is unreachable.
func NewClient(ctx context.Context, url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("cache: parse url: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("cache: ping: %w", err)
	}

	return client, nil
}
