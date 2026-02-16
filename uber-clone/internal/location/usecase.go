package location

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
)

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) GetNearbyDrivers(ctx context.Context, location *domain.Location) ([]*domain.Driver, error) {
	return nil, nil
}

func (u *UseCase) GetDriverLatestLocation(ctx context.Context, driverID string) (*domain.Location, error) {
	return nil, nil
}

func (u *UseCase) UpdateDriverLatestLocation(ctx context.Context, location *domain.Location) error {
	return nil
}
