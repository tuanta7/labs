package main

import (
	"time"
)

type SagaState string

const (
	SagaStarted         SagaState = "STARTED"
	SagaReservePending  SagaState = "RESERVE_PENDING"
	SagaReserveOK       SagaState = "RESERVE_OK"
	SagaPaymentPending  SagaState = "PAYMENT_PENDING"
	SagaPaymentOK       SagaState = "PAYMENT_OK"
	SagaChargingPending SagaState = "CHARGING_PENDING"
	SagaCompleted       SagaState = "COMPLETED"
	SagaCompensating    SagaState = "COMPENSATING"
	SagaFailed          SagaState = "FAILED"
)

// Saga represents the orchestrated saga state persisted in MongoDB.
type Saga struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	StationID string    `json:"station_id" bson:"station_id"`
	Amount    float64   `json:"amount" bson:"amount"`
	Status    SagaState `json:"status" bson:"status"`
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
var StepToPendingStatus = map[StepType]SagaState{
	StepReserve:  SagaReservePending,
	StepPayment:  SagaPaymentPending,
	StepCharging: SagaChargingPending,
}

// StepToOKStatus maps a step to its completed saga status.
var StepToOKStatus = map[StepType]SagaState{
	StepReserve:  SagaReserveOK,
	StepPayment:  SagaPaymentOK,
	StepCharging: SagaCompleted,
}

// StepToQueue maps a step to its RabbitMQ queue/routing key.
var StepToQueue = map[StepType]string{
	StepReserve:  QueueReservation,
	StepPayment:  QueuePayment,
	StepCharging: QueueCharging,
}
