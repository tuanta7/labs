package kafka

import (
	"context"
	"encoding/json"
	"kafka-lab/internal/domain"
	"time"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	logger *zerolog.Logger
}

func NewProducer(seeds []string) (*Producer, error) {
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

	return &Producer{
		client: c,
		logger: zerolog.DefaultContextLogger,
	}, nil
}

func (p *Producer) Close() {
	p.client.Close()
}

func (p *Producer) ProduceSync(
	ctx context.Context,
	topic string,
	key []byte,
	location *domain.Location,
) error {
	value, err := json.Marshal(location)
	if err != nil {
		return err
	}

	record, err := p.client.ProduceSync(ctx, &kgo.Record{
		Topic:     topic,
		Key:       key,
		Value:     value,
		Timestamp: time.Now().UTC(),
	}).First()

	if err != nil {
		p.logger.Error().
			Err(err).
			Str("topic", topic).
			Str("key", string(key)).
			Msg("failed to produce record")
		return err
	}

	p.logger.Info().
		Str("topic", topic).
		Str("key", string(key)).
		Int64("offset", record.Offset).
		Msg("message produced to kafka successfully")

	return nil
}
