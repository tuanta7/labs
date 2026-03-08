package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ChargingRequest is the inbound request to start an EV charging saga.
type ChargingRequest struct {
	UserID    string  `json:"user_id"`
	StationID string  `json:"station_id"`
	Amount    float64 `json:"amount"`
}

// StartSaga creates a new saga, persists it, and publishes the first command (Reserve).
func StartSaga(ctx context.Context, ch *amqp.Channel, store *SagaStore, req ChargingRequest) (*Saga, error) {
	saga := &Saga{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		StationID: req.StationID,
		Amount:    req.Amount,
		Status:    SagaStarted,
	}

	if err := store.CreateSaga(ctx, saga); err != nil {
		return nil, fmt.Errorf("create saga: %w", err)
	}

	// Transition to RESERVE_PENDING.
	if err := store.UpdateStatus(ctx, saga.ID, SagaReservePending, ""); err != nil {
		return nil, fmt.Errorf("update saga to reserve pending: %w", err)
	}
	saga.Status = SagaReservePending

	// Publish the first command: Reserve the charging station.
	cmd := SagaCommand{
		SagaID:    saga.ID,
		Step:      StepReserve,
		Action:    ActionExecute,
		UserID:    saga.UserID,
		StationID: saga.StationID,
		Amount:    saga.Amount,
	}
	if err := PublishCommand(ctx, ch, cmd); err != nil {
		return nil, fmt.Errorf("publish reserve command: %w", err)
	}

	log.Printf("[orchestrator] saga %s STARTED — published RESERVE command", saga.ID)
	return saga, nil
}
