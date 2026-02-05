package redis

import "github.com/redis/go-redis/extra/redisotel/v9"

type Option func(*Client) error

func WithTraces() Option {
	return func(c *Client) error {
		return redisotel.InstrumentMetrics(c.redisClient)
	}
}

func WithMetrics() Option {
	return func(c *Client) error {
		return redisotel.InstrumentTracing(c.redisClient)
	}
}
