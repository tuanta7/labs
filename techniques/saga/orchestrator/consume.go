package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeCommands  = "saga.commands"
	ExchangeResponses = "saga.responses"

	QueueReservation = "reservation"
	QueuePayment     = "payment"
	QueueCharging    = "charging"
	QueueResponse    = "response"

	PrefetchCount = 1
)

func SetupTopology(ch *amqp.Channel) error {
	// Declare the commands exchange (topic) used by the orchestrator to publish commands.
	if err := ch.ExchangeDeclare(ExchangeCommands, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare commands exchange: %w", err)
	}

	// Declare command queues and bind them to the commands exchange.
	commandQueues := []string{QueueReservation, QueuePayment, QueueCharging}
	for _, q := range commandQueues {
		if _, err := ch.QueueDeclare(q, true, false, false, false, nil); err != nil {
			return fmt.Errorf("declare queue %s: %w", q, err)
		}
		if err := ch.QueueBind(q, q, ExchangeCommands, false, nil); err != nil {
			return fmt.Errorf("bind queue %s: %w", q, err)
		}
	}

	// Declare the responses exchange (direct) used by services to reply.
	if err := ch.ExchangeDeclare(ExchangeResponses, amqp.ExchangeDirect, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare responses exchange: %w", err)
	}
	// Declare the response queue and bind it to the responses exchange.
	if _, err := ch.QueueDeclare(QueueResponse, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare response queue: %w", err)
	}
	if err := ch.QueueBind(QueueResponse, QueueResponse, ExchangeResponses, false, nil); err != nil {
		return fmt.Errorf("bind response queue: %w", err)
	}

	return nil
}

// PublishCommand sends a SagaCommand to the appropriate service queue.
func PublishCommand(ctx context.Context, ch *amqp.Channel, cmd SagaCommand) error {
	routingKey, ok := StepToQueue[cmd.Step]
	if !ok {
		return fmt.Errorf("unknown step: %s", cmd.Step)
	}

	body, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshal command: %w", err)
	}

	return ch.PublishWithContext(ctx, ExchangeCommands, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

// ConsumeResponses listens to the response queue and drives the saga state machine.
func ConsumeResponses(ctx context.Context, ch *amqp.Channel, store *SagaStore) error {
	// Set QoS: process one message at a time per consumer.
	if err := ch.Qos(PrefetchCount, 0, false); err != nil {
		return fmt.Errorf("set QoS: %w", err)
	}

	deliveries, err := ch.Consume(
		QueueResponse,
		"orchestrator-consumer",
		false, // manual acknowledgment
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume response queue: %w", err)
	}

	log.Println("[orchestrator] listening for saga responses...")
	for {
		select {
		case <-ctx.Done():
			log.Println("[orchestrator] context cancelled, stopping consumer")
			return ctx.Err()
		case d, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("response channel closed")
			}

			if err = handleResponse(ctx, ch, store, d); err != nil {
				log.Printf("[orchestrator] error handling response: %v, nacking", err)
				// Nack and requeue so we can retry.
				_ = d.Nack(false, true)
			}
		}
	}
}

// handleResponse processes a single SagaResponse delivery.
func handleResponse(ctx context.Context, ch *amqp.Channel, store *SagaStore, d amqp.Delivery) error {
	var resp SagaResponse
	if err := json.Unmarshal(d.Body, &resp); err != nil {
		// Bad message format — ack to discard (no point retrying).
		log.Printf("[orchestrator] invalid message, discarding: %v", err)
		return d.Ack(false)
	}

	log.Printf("[orchestrator] received response: saga=%s step=%s success=%v err=%s", resp.SagaID, resp.Step, resp.Success, resp.Error)

	saga, err := store.GetSaga(ctx, resp.SagaID)
	if err != nil {
		return fmt.Errorf("get saga %s: %w", resp.SagaID, err)
	}

	if resp.Success {
		return handleStepSuccess(ctx, ch, store, saga, resp, d)
	}

	return handleStepFailure(ctx, ch, store, saga, resp, d)
}

// handleStepSuccess advances the saga to the next step or marks it as completed.
func handleStepSuccess(ctx context.Context, ch *amqp.Channel, store *SagaStore, saga *Saga, resp SagaResponse, d amqp.Delivery) error {
	// Mark current step as OK.
	okStatus, exists := StepToOKStatus[resp.Step]
	if !exists {
		return fmt.Errorf("unknown step: %s", resp.Step)
	}
	if err := store.UpdateStatus(ctx, saga.ID, okStatus, ""); err != nil {
		return fmt.Errorf("update saga status: %w", err)
	}

	// Find the next step.
	nextStep, hasNext := nextSagaStep(resp.Step)
	if !hasNext {
		// All steps completed!
		log.Printf("[orchestrator] saga %s COMPLETED", saga.ID)
		return d.Ack(false)
	}

	// Transition to pending for the next step.
	pendingStatus := StepToPendingStatus[nextStep]
	if err := store.UpdateStatus(ctx, saga.ID, pendingStatus, ""); err != nil {
		return fmt.Errorf("update saga to pending: %w", err)
	}

	// Publish the next command.
	cmd := SagaCommand{
		SagaID:    saga.ID,
		Step:      nextStep,
		Action:    ActionExecute,
		UserID:    saga.UserID,
		StationID: saga.StationID,
		Amount:    saga.Amount,
	}
	if err := PublishCommand(ctx, ch, cmd); err != nil {
		return fmt.Errorf("publish next command: %w", err)
	}

	log.Printf("[orchestrator] saga %s: advanced to step %s", saga.ID, nextStep)
	return d.Ack(false)
}

// handleStepFailure initiates compensation for all previously completed steps.
func handleStepFailure(ctx context.Context, ch *amqp.Channel, store *SagaStore, saga *Saga, resp SagaResponse, d amqp.Delivery) error {
	log.Printf("[orchestrator] saga %s: step %s FAILED: %s — starting compensation", saga.ID, resp.Step, resp.Error)

	if err := store.UpdateStatus(ctx, saga.ID, SagaCompensating, resp.Error); err != nil {
		return fmt.Errorf("update saga to compensating: %w", err)
	}

	// Find all steps that completed before the failed step, and compensate in reverse.
	failedIdx := stepIndex(resp.Step)
	for i := failedIdx - 1; i >= 0; i-- {
		step := SagaSteps[i]
		cmd := SagaCommand{
			SagaID:    saga.ID,
			Step:      step,
			Action:    ActionRollback,
			UserID:    saga.UserID,
			StationID: saga.StationID,
			Amount:    saga.Amount,
		}
		if err := PublishCommand(ctx, ch, cmd); err != nil {
			return fmt.Errorf("publish compensate command for step %s: %w", step, err)
		}
		log.Printf("[orchestrator] saga %s: published ROLLBACK for step %s", saga.ID, step)
	}

	if err := store.UpdateStatus(ctx, saga.ID, SagaFailed, resp.Error); err != nil {
		return fmt.Errorf("update saga to failed: %w", err)
	}

	return d.Ack(false)
}

// nextSagaStep returns the step after the given step, if any.
func nextSagaStep(current StepType) (StepType, bool) {
	idx := stepIndex(current)
	if idx < 0 || idx >= len(SagaSteps)-1 {
		return "", false
	}
	return SagaSteps[idx+1], true
}

// stepIndex returns the index of a step in SagaSteps, or -1.
func stepIndex(step StepType) int {
	for i, s := range SagaSteps {
		if s == step {
			return i
		}
	}
	return -1
}
