package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type LocationKafkaProducer struct {
	client *kgo.Client
}

func NewProducer(seeds []string) (*LocationKafkaProducer, error) {
	c, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.AllowAutoTopicCreation(),
		kgo.DialTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = c.Ping(ctx); err != nil {
		return nil, err
	}

	return &LocationKafkaProducer{
		client: c,
	}, nil
}

func (p *LocationKafkaProducer) Close() {
	p.client.Close()
}

func (p *LocationKafkaProducer) ProduceSync(
	ctx context.Context,
	topic string,
	key []byte,
	location *Location,
) error {
	value, err := json.Marshal(location)
	if err != nil {
		return err
	}

	return p.client.ProduceSync(ctx, &kgo.Record{
		Topic:     topic,
		Key:       key,
		Value:     value,
		Timestamp: time.Now().UTC(),
	}).FirstErr()
}
