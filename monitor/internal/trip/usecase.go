package trip

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tuanta7/monitor/pkg/monitor"
	"go.uber.org/zap"
)

type UseCase struct {
	repo   *Repository
	logger *monitor.Logger
	meter  *monitor.Meter
	tracer *monitor.Tracer
}

func NewUseCase(
	repo *Repository,
	logger *monitor.Logger,
	meter *monitor.Meter,
	tracer *monitor.Tracer,
) *UseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
		meter:  meter,
		tracer: tracer,
	}
}

func (u *UseCase) CreateTrip(ctx context.Context, trip *Trip) error {
	trip.Status = StatusPending
	trip.ID = uuid.NewString()
	trip.CreatedAt = time.Now().UTC()
	trip.UpdatedAt = trip.CreatedAt

	ctx, span := u.tracer.StartNewSpan(ctx, "create-trip")
	defer span.End()

	err := u.repo.CreateTrip(ctx, trip)
	if err != nil {
		span.RecordError(err)
		u.logger.Error("failed to create trip",
			zap.Error(err),
			zap.String("tripID", trip.ID),
		)
		return err
	}

	return nil
}

func (u *UseCase) AcceptTrip(ctx context.Context, tripID, driverID string) error {
	err := u.repo.AcceptTrip(ctx, &Trip{
		ID:       tripID,
		DriverID: driverID,
	})
	if err != nil {
		u.logger.Error("failed to get driver",
			zap.Error(err),
			zap.String("driverID", driverID),
		)
		return err
	}

	return nil
}
