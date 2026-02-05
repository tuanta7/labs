package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
}

func NewConsumer(ctx context.Context, seeds []string, topic, group string) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client: client,
	}, nil
}

func (c *Consumer) Close() {
	c.client.Close()
}

type ConsumerHandler func(ctx context.Context, key, value []byte)

func (c *Consumer) Consume(ctx context.Context, handler ConsumerHandler) error {
	for {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			continue
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, record := range p.Records {
				handler(ctx, record.Key, record.Value)
			}
		})
	}
}
