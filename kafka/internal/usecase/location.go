package usecase

import (
	"context"
	"encoding/json"
	"kafka-lab/internal/domain"
	"kafka-lab/internal/repository"
	"time"

	"github.com/rs/zerolog"
)

type Publisher interface {
	ProduceSync(ctx context.Context, topic string, key []byte, location *domain.Location) error
}

type Options func(*LocationUC)

func WithPublisher(p Publisher, topic string) Options {
	return func(l *LocationUC) {
		l.publisher = p
		l.topic = topic
	}
}

func WithRepository(r *repository.LocationRepository) Options {
	return func(l *LocationUC) {
		l.locationRepo = r
	}
}

type LocationUC struct {
	topic        string
	publisher    Publisher
	locationRepo *repository.LocationRepository
	logger       *zerolog.Logger
}

func NewLocationUC(
	logger *zerolog.Logger,
	options ...Options,
) *LocationUC {
	if logger == nil {
		logger = zerolog.DefaultContextLogger
	}

	uc := &LocationUC{
		logger: logger,
	}

	for _, opt := range options {
		opt(uc)
	}

	return uc
}

func (l *LocationUC) ProcessLocation(ctx context.Context, location *domain.Location) {
	key := []byte(location.UserID)

	if err := l.publisher.ProduceSync(ctx, l.topic, key, location); err != nil {
		l.logger.Error().
			Err(err).
			Str("location_id", location.ID).
			Str("topic", l.topic).
			Msg("failed to produce location to kafka")
		return
	}

	l.logger.Info().
		Str("location_id", location.ID).
		Str("user_id", location.UserID).
		Str("topic", l.topic).
		Msg("location successfully produced to kafka")
}

func (l *LocationUC) SaveLocation(ctx context.Context, value []byte) error {
	var location domain.Location
	if err := json.Unmarshal(value, &location); err != nil {
		return err
	}

	if location.Timestamp == 0 {
		location.Timestamp = time.Now().UTC().UnixMilli()
	}

	if err := l.locationRepo.SaveLocation(ctx, &location); err != nil {
		return err
	}

	return nil
}
