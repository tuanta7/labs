package trip

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tuanta7/monitor/pkg/monitor"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
)

type UseCase struct {
	repo   *Repository
	client *http.Client
	logger *monitor.Logger
	meter  *monitor.Meter
	tracer *monitor.Tracer
}

func NewUseCase(
	repo *Repository,
	client *http.Client,
	logger *monitor.Logger,
	meter *monitor.Meter,
	tracer *monitor.Tracer,
) *UseCase {
	return &UseCase{
		repo:   repo,
		client: client,
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

	ctx1, span1 := u.tracer.StartNewSpan(ctx, "trip.create")
	span1.SetAttributes(
		semconv.DBMongoDBCollection(TripsCollection),
		semconv.DBOperation("insert"),
	)

	err := u.repo.CreateTrip(ctx1, trip)
	if err != nil {
		span1.RecordError(err)
		u.logger.Error("failed to create trip",
			zap.Error(err),
			zap.String("tripID", trip.ID),
		)
		return err
	}
	span1.End()

	ctx2, span2 := u.tracer.StartNewSpan(ctx, "trip.create.notiy")
	defer span2.End()

	req, err := http.NewRequestWithContext(ctx2, http.MethodPost, "http://localhost:13072/notify", nil)
	if err != nil {
		span2.RecordError(err)
		u.logger.Error("failed to create notification request", zap.Error(err))
		return err
	}
	otel.GetTextMapPropagator().Inject(ctx2, propagation.HeaderCarrier(req.Header))

	response, err := u.client.Do(req)
	if err != nil {
		span2.RecordError(err)
		u.logger.Error("failed to send notification", zap.Error(err))
		return err
	}
	defer response.Body.Close()

	span2.SetAttributes(
		semconv.HTTPStatusCode(response.StatusCode),
		semconv.HTTPMethod(req.Method),
		semconv.HTTPTarget(req.URL.String()),
	)

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
