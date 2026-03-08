package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeCommands  = "saga.commands"
	ExchangeResponses = "saga.responses"
	QueuePayment      = "payment"
	QueueResponse     = "response"
	PrefetchCount     = 1
)

// SagaCommand mirrors the orchestrator's SagaCommand.
type SagaCommand struct {
	SagaID    string  `json:"saga_id"`
	Step      string  `json:"step"`
	Action    string  `json:"action"`
	UserID    string  `json:"user_id"`
	StationID string  `json:"station_id"`
	Amount    float64 `json:"amount"`
}

// SagaResponse is sent back to the orchestrator.
type SagaResponse struct {
	SagaID  string `json:"saga_id"`
	Step    string `json:"step"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rabbitURL := envOrDefault("RABBITMQ_URL", "amqp://rabbitmq:password@localhost:5672/")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("[payment] connect to rabbitmq: %v", err)
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[payment] open channel: %v", err)
	}
	defer func() { _ = ch.Close() }()

	// Set QoS — process one message at a time.
	if err := ch.Qos(PrefetchCount, 0, false); err != nil {
		log.Fatalf("[payment] set QoS: %v", err)
	}

	// Ensure topology exists (idempotent declarations).
	if err := ch.ExchangeDeclare(ExchangeCommands, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		log.Fatalf("[payment] declare commands exchange: %v", err)
	}
	if err := ch.ExchangeDeclare(ExchangeResponses, amqp.ExchangeDirect, true, false, false, false, nil); err != nil {
		log.Fatalf("[payment] declare responses exchange: %v", err)
	}
	if _, err := ch.QueueDeclare(QueuePayment, true, false, false, false, nil); err != nil {
		log.Fatalf("[payment] declare queue: %v", err)
	}
	if err := ch.QueueBind(QueuePayment, QueuePayment, ExchangeCommands, false, nil); err != nil {
		log.Fatalf("[payment] bind queue: %v", err)
	}
	if _, err := ch.QueueDeclare(QueueResponse, true, false, false, false, nil); err != nil {
		log.Fatalf("[payment] declare response queue: %v", err)
	}
	if err := ch.QueueBind(QueueResponse, QueueResponse, ExchangeResponses, false, nil); err != nil {
		log.Fatalf("[payment] bind response queue: %v", err)
	}

	deliveries, err := ch.Consume(
		QueuePayment,
		"payment-consumer",
		false, // autoAck = false
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[payment] consume: %v", err)
	}

	log.Println("[payment] service started, waiting for commands...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[payment] shutting down")
			return
		case d, ok := <-deliveries:
			if !ok {
				log.Println("[payment] channel closed")
				return
			}
			handleDelivery(ctx, ch, d)
		}
	}
}

func handleDelivery(ctx context.Context, ch *amqp.Channel, d amqp.Delivery) {
	var cmd SagaCommand
	if err := json.Unmarshal(d.Body, &cmd); err != nil {
		log.Printf("[payment] invalid message, discarding: %v", err)
		_ = d.Ack(false)
		return
	}

	log.Printf("[payment] received: saga=%s action=%s user=%s amount=%.2f",
		cmd.SagaID, cmd.Action, cmd.UserID, cmd.Amount)

	var resp SagaResponse
	resp.SagaID = cmd.SagaID
	resp.Step = cmd.Step

	switch cmd.Action {
	case "EXECUTE":
		resp = authorizePayment(cmd)
	case "ROLLBACK":
		resp = refundPayment(cmd)
	default:
		log.Printf("[payment] unknown action: %s", cmd.Action)
		resp.Success = false
		resp.Error = fmt.Sprintf("unknown action: %s", cmd.Action)
	}

	if err := publishResponse(ctx, ch, resp); err != nil {
		log.Printf("[payment] failed to publish response: %v — nacking", err)
		_ = d.Nack(false, true)
		return
	}

	_ = d.Ack(false)
}

func authorizePayment(cmd SagaCommand) SagaResponse {
	resp := SagaResponse{SagaID: cmd.SagaID, Step: cmd.Step}

	// Simulate payment authorization.
	// 85% success rate for demonstration purposes.
	if rand.Intn(100) < 85 {
		log.Printf("[payment] saga=%s: payment of %.2f authorized for user %s",
			cmd.SagaID, cmd.Amount, cmd.UserID)
		resp.Success = true
	} else {
		log.Printf("[payment] saga=%s: payment of %.2f DECLINED for user %s (simulated)",
			cmd.SagaID, cmd.Amount, cmd.UserID)
		resp.Success = false
		resp.Error = "insufficient funds"
	}
	return resp
}

func refundPayment(cmd SagaCommand) SagaResponse {
	resp := SagaResponse{SagaID: cmd.SagaID, Step: cmd.Step}

	// Simulate refund — always succeeds.
	log.Printf("[payment] saga=%s: payment of %.2f REFUNDED for user %s (compensating)",
		cmd.SagaID, cmd.Amount, cmd.UserID)
	resp.Success = true
	return resp
}

func publishResponse(ctx context.Context, ch *amqp.Channel, resp SagaResponse) error {
	body, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	return ch.PublishWithContext(ctx, ExchangeResponses, QueueResponse, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
