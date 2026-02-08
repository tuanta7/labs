package ingestion

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/kafka"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
)

type UseCase struct {
	logger    *zapx.Logger
	publisher *kafka.Publisher
}

func NewUseCase(logger *zapx.Logger, publisher *kafka.Publisher) *UseCase {
	return &UseCase{
		logger:    logger,
		publisher: publisher,
	}
}

func (u *UseCase) PublishLocation(ctx context.Context, location *domain.DriverLocationMessage) {
	key := []byte(location.DriverID)

	u.logger.Debug("publishing driver location",
		zap.String("topic", domain.DriverLocationTopic),
		zap.String("key", location.DriverID),
	)

	u.publisher.Publish(ctx, domain.DriverLocationTopic, key, []byte("msg"), nil)
}
