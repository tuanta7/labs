package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
	logger *zerolog.Logger
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
		logger: zerolog.DefaultContextLogger,
	}, nil
}

func (c *Consumer) Close() {
	c.client.Close()
}

func (c *Consumer) Consume(ctx context.Context, handler func(key, value []byte) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fetches := c.client.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		for _, fetchErr := range errs {
			if errors.Is(fetchErr.Err, context.Canceled) || errors.Is(fetchErr.Err, context.DeadlineExceeded) {
				c.logger.Info().Err(fetchErr.Err).Msg("kafka fetch canceled")
				return fetchErr.Err
			}

			c.logger.Error().Err(fetchErr.Err).Msg("kafka fetch failed")
		}

		return fmt.Errorf("kafka fetch failed: %w", errs[0].Err)
	}

	var handlerErr error
	fetches.EachPartition(func(p kgo.FetchTopicPartition) {
		for _, record := range p.Records {
			if handlerErr != nil {
				return
			}

			if err := handler(record.Key, record.Value); err != nil {
				handlerErr = err
				return
			}
		}
	})

	if handlerErr != nil {
		return handlerErr
	}

	return nil
}
