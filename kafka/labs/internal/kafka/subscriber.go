package kafka

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type MessageHandler func(ctx context.Context, msg *kgo.Record) error

type Consumer struct {
	client *kgo.Client
}

func NewConsumer(seeds []string, topic, group string) (*Consumer, error) {
	c, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = c.Ping(ctx); err != nil {
		return nil, err
	}

	return &Consumer{
		client: c,
	}, nil
}

func (c *Consumer) Close() {
	c.client.Close()
}

func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) (err error) {
	fetches := c.client.PollFetches(ctx)

	fetches.EachRecord(func(r *kgo.Record) {
		if err != nil {
			return
		}

		handlerErr := handler(ctx, r)
		if err != nil {
			// log / retry / DLQ
			err = handlerErr
		}

		_ = c.client.CommitRecords(ctx, r)
	})

	return err
}
