package redis

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	MGet(ctx context.Context, keys ...string) ([]any, error)
	MSet(ctx context.Context, values ...any) error
	Exists(ctx context.Context, keys ...string) (int64, error)
}

func (c *Client) Exists(ctx context.Context, key ...string) (int64, error) {
	result := c.redisClient.Exists(ctx, key...)
	if err := result.Err(); err != nil {
		return 0, err
	}

	return result.Val(), nil
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	result := c.redisClient.Get(ctx, key)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Bytes()
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.redisClient.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.redisClient.Del(ctx, keys...).Err()
}

func (c *Client) MGet(ctx context.Context, keys ...string) ([]any, error) {
	result := c.redisClient.MGet(ctx, keys...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Result()
}

func (c *Client) MSet(ctx context.Context, pairs ...any) error {
	return c.redisClient.MSet(ctx, pairs...).Err()
}
