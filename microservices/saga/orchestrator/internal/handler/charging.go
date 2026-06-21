package handler

import (
	"context"
	"fmt"
	"log"
	"orchestrator/internal/saga"

	"github.com/google/uuid"
)

// ChargingRequest is the inbound request to start an EV charging saga.
type ChargingRequest struct {
	RequestID string  `json:"request_id"`
	StationID string  `json:"station_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
}

// StartChargingSaga creates a new saga, persists it, and publishes the first command (Reserve).
func (sv *Server) StartChargingSaga(ctx context.Context, req ChargingRequest) (*saga.Saga, error) {
	s := &saga.Saga{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		StationID: req.StationID,
		Amount:    req.Amount,
		Status:    saga.StateStarted,
	}

	if err := sv.store.CreateSaga(ctx, s); err != nil {
		return nil, fmt.Errorf("create saga: %w", err)
	}

	if err := sv.store.UpdateStatus(ctx, s.ID, saga.StateReservePending, ""); err != nil {
		return nil, fmt.Errorf("update saga to reserve pending: %w", err)
	}
	s.Status = saga.StateReservePending

	cmd := saga.Command{
		SagaID:    s.ID,
		Step:      saga.StepReserve,
		Action:    saga.ActionExecute,
		UserID:    s.UserID,
		StationID: s.StationID,
		Amount:    s.Amount,
	}
	if err := PublishCommand(ctx, sv.ch, cmd); err != nil {
		return nil, fmt.Errorf("publish reserve command: %w", err)
	}

	log.Printf("[orchestrator] saga %s STARTED — published RESERVE command", s.ID)
	return s, nil
}
