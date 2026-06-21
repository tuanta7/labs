package saga

import (
	"time"
)

type State string

const (
	StateStarted         State = "STARTED"
	StateReservePending  State = "RESERVE_PENDING"
	StateReserveOK       State = "RESERVE_OK"
	StatePaymentPending  State = "PAYMENT_PENDING"
	StatePaymentOK       State = "PAYMENT_OK"
	StateChargingPending State = "CHARGING_PENDING"
	StateCompleted       State = "COMPLETED"
	StateCompensating    State = "COMPENSATING"
	StateFailed          State = "FAILED"
)

// Saga represents the orchestrated saga state persisted in MongoDB.
type Saga struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	StationID string    `json:"station_id" bson:"station_id"`
	Amount    float64   `json:"amount" bson:"amount"`
	Status    State     `json:"status" bson:"status"`
	Error     string    `json:"error,omitempty" bson:"error,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// SagaSteps defines the ordered forward steps of the saga.
var SagaSteps = []StepType{
	StepReserve,
	StepPayment,
	StepCharging,
}

// StepToPendingStatus maps a step to its pending saga status.
var StepToPendingStatus = map[StepType]State{
	StepReserve:  StateReservePending,
	StepPayment:  StatePaymentPending,
	StepCharging: StateChargingPending,
}

// StepToOKStatus maps a step to its completed saga status.
var StepToOKStatus = map[StepType]State{
	StepReserve:  StateReserveOK,
	StepPayment:  StatePaymentOK,
	StepCharging: StateCompleted,
}
