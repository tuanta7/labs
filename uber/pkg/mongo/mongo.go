package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
)

type Config struct {
	URI            string
	Database       string
	MaxPoolSize    uint64
	MinPoolSize    uint64
	ConnectTimeout time.Duration
	QueryTimeout   time.Duration
	Monitor        bool
}

type Client struct {
	client   *mongo.Client
	database *mongo.Database
	timeout  time.Duration
}

func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	opts := options.Client().
		ApplyURI(cfg.URI).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetConnectTimeout(cfg.ConnectTimeout)

	if cfg.Monitor {
		opts.SetMonitor(otelmongo.NewMonitor())
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if pingErr := client.Ping(ctx, readpref.Primary()); pingErr != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", pingErr)
	}

	timeout := 10 * time.Second
	if cfg.QueryTimeout > 0 {
		timeout = cfg.QueryTimeout
	}

	return &Client{
		client:   client,
		database: client.Database(cfg.Database),
		timeout:  timeout,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *Client) Collection(name string) *CollectionClient {
	return &CollectionClient{
		timeout:    c.timeout,
		collection: c.database.Collection(name),
	}
}
