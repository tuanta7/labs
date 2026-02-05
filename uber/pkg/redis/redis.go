package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type Config struct {
	MasterName    string
	Username      string
	Password      string
	SentinelAddrs []string
}

type Client struct {
	redisClient *goredis.Client
}

func NewFailoverClient(ctx context.Context, config *Config, opts ...Option) (*Client, error) {
	client := &Client{}

	c := goredis.NewFailoverClient(&goredis.FailoverOptions{
		MasterName:    config.MasterName,
		SentinelAddrs: config.SentinelAddrs,
		Username:      config.Username,
		Password:      config.Password,
	})
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	client.redisClient = c
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func NewClient(ctx context.Context, config *Config, opts ...Option) (*Client, error) {
	client := &Client{}

	c := goredis.NewClient(&goredis.Options{
		Addr:     "",
		Username: config.Username,
		Password: config.Password,
	})
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	client.redisClient = c
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) Close() error {
	return c.redisClient.Close()
}
