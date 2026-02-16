package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Collection interface {
	InsertOne(ctx context.Context, document any) error
	FindOne(ctx context.Context, filter any, result any) error
	UpdateOne(ctx context.Context, filter any, update any) error
}

type CollectionClient struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func (c *CollectionClient) FindOne(ctx context.Context, filter any, result any) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.collection.FindOne(timeoutCtx, filter).Decode(result)
}

func (c *CollectionClient) InsertOne(ctx context.Context, document any) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.collection.InsertOne(timeoutCtx, document)
	return err
}

func (c *CollectionClient) UpdateOne(ctx context.Context, filter any, update any) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.collection.UpdateOne(timeoutCtx, filter, update)
	return err
}
