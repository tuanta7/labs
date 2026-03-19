package usecase

import (
	"context"
	"encoding/json"
	"kafka-lab/internal/domain"
	"kafka-lab/internal/repository"
	"time"

	"github.com/rs/zerolog"
)

type Consumer interface {
	Consume(ctx context.Context, handler func(key, value []byte) error) error
}

type ConsumerUC struct {
	locationRepo *repository.LocationRepository
	consumer     Consumer
	topic        string
	groupID      string
	logger       *zerolog.Logger
}

func NewConsumerUseCase(
	consumer Consumer,
	locationRepo *repository.LocationRepository,
	topic string,
	groupID string,
	logger *zerolog.Logger,
) *ConsumerUC {
	if logger == nil {
		logger = zerolog.DefaultContextLogger
	}

	return &ConsumerUC{
		consumer:     consumer,
		locationRepo: locationRepo,
		topic:        topic,
		groupID:      groupID,
		logger:       logger,
	}
}

func (c *ConsumerUC) ConsumeLocationTopic(ctx context.Context) error {
	for {
		err := c.consumer.Consume(ctx, func(key, value []byte) error {
			var location domain.Location
			if err := json.Unmarshal(value, &location); err != nil {
				return err
			}

			if location.Timestamp == 0 {
				location.Timestamp = time.Now().UTC().UnixMilli()
			}

			if err := c.locationRepo.SaveLocation(ctx, &location); err != nil {
				return err
			}

			c.logger.Info().
				Bytes("key", key).
				Str("location_id", location.ID).
				Str("user_id", location.UserID).
				Str("topic", c.topic).
				Str("group_id", c.groupID).
				Msg("consumed and persisted location")

			return nil
		})

		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			c.logger.Error().Err(err).
				Str("topic", c.topic).
				Str("group_id", c.groupID).
				Msg("consumer loop failed, retrying")

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
			}

			continue
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
