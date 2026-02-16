package driver

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
)

type UseCase struct {
	logger *zapx.Logger
	repo   *Repository
}

func NewUseCase(logger *zapx.Logger, repo *Repository) *UseCase {
	return &UseCase{
		logger: logger,
		repo:   repo,
	}
}

func (u *UseCase) GetDriverByID(ctx context.Context, driverID string) (*domain.Driver, error) {
	driver, err := u.repo.GetDriverByID(ctx, driverID)
	if err != nil {
		u.logger.Error("failed to get driver",
			zap.Error(err),
			zap.String("driverID", driverID),
		)
		return nil, err
	}

	return driver, nil
}
