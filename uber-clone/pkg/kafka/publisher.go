package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Publisher struct {
	client *kgo.Client
}

func NewPublisher(seeds []string) (*Publisher, error) {
	c, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		client: c,
	}, nil
}

func (p *Publisher) Close() {
	p.client.Close()
}

func (p *Publisher) Publish(ctx context.Context, topic string, key, value []byte, cb func(err error)) {
	p.client.Produce(ctx, &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	}, func(record *kgo.Record, err error) {
		cb(err)
	})
}

func (p *Publisher) PublishSync(ctx context.Context, topic string, key, value []byte) error {
	return p.client.ProduceSync(ctx, &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	}).FirstErr()
}
