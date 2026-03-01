package main

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumer struct {
	client *kgo.Client
}

func NewConsumer(ctx context.Context, seeds []string, topic, group string) (*KafkaConsumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		client: client,
	}, nil
}

func (c *KafkaConsumer) Close() {
	c.client.Close()
}

func (c *KafkaConsumer) ConsumeLocationMessage(ctx context.Context) error {
	for {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			continue
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, record := range p.Records {
				fmt.Printf("record: %s\n", record.Value)
			}
		})
	}
}
