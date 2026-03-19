package usecase

import (
	"context"
	"kafka-lab/internal/domain"

	"github.com/rs/zerolog"
)

type Publisher interface {
	ProduceSync(ctx context.Context, topic string, key []byte, location *domain.Location) error
}

type PublisherUC struct {
	publisher Publisher
	topic     string
	logger    *zerolog.Logger
}

func NewPublisherUseCase(producer Publisher, topic string, logger *zerolog.Logger) *PublisherUC {
	return &PublisherUC{
		publisher: producer,
		topic:     topic,
		logger:    logger,
	}
}

func (p *PublisherUC) PublishLocationMessage(ctx context.Context, location *domain.Location) {
	key := []byte(location.UserID)

	if err := p.publisher.ProduceSync(ctx, p.topic, key, location); err != nil {
		p.logger.Error().
			Err(err).
			Str("location_id", location.ID).
			Str("topic", p.topic).
			Msg("failed to produce location to kafka")
		return
	}

	p.logger.Info().
		Str("location_id", location.ID).
		Str("user_id", location.UserID).
		Str("topic", p.topic).
		Msg("location successfully produced to kafka")
}
