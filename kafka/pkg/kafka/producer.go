package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
}

func NewProducer(seeds []string) (*Producer, error) {
	c, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.AllowAutoTopicCreation(), //
	)
	if err != nil {
		return nil, err
	}

	return &Producer{
		client: c,
	}, nil
}

func (p *Producer) Close() {
	p.client.Close()
}

func (p *Producer) Produce(ctx context.Context, topic string, key, value []byte) {
	p.client.Produce(ctx, &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	}, func(record *kgo.Record, err error) {
		if err != nil {

		}
	})
}
