package main

import (
	"context"

	"github.com/rs/zerolog"
)

type LocationProducer interface {
	ProduceSync(ctx context.Context, topic string, key []byte, location *Location) error
}

type LocationUseCase struct {
	producer LocationProducer
	topic    string
	logger   *zerolog.Logger
}

func NewUseCase(producer LocationProducer, topic string, logger *zerolog.Logger) *LocationUseCase {
	return &LocationUseCase{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}
}

func (p *LocationUseCase) ProduceLocation(ctx context.Context, location *Location) {
	key := []byte(location.UserID)

	if err := p.producer.ProduceSync(ctx, p.topic, key, location); err != nil {
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
