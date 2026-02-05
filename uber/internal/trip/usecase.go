package trip

import (
	"context"

	"github.com/google/uuid"
	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
)

type UseCase struct {
	db     Repository
	logger *zapx.Logger
}

func NewUseCase(db Repository, logger *zapx.Logger) *UseCase {
	return &UseCase{db: db, logger: logger}
}

func (u *UseCase) CreateTrip(ctx context.Context, trip *domain.Trip) error {
	if trip.ID == "" {
		trip.ID = uuid.NewString()
	}

	return nil
}

func (u *UseCase) GetTripByID(ctx context.Context, tripID string) (*domain.Trip, error) {
	return nil, nil
}

func (u *UseCase) InsertLocations(ctx context.Context, locations []*domain.Location) error {
	return nil
}
