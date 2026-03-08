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
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeCommands  = "saga.commands"
	ExchangeResponses = "saga.responses"
	QueueCharging     = "charging"
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
		log.Fatalf("[charging] connect to rabbitmq: %v", err)
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[charging] open channel: %v", err)
	}
	defer func() { _ = ch.Close() }()

	// Set QoS — process one message at a time.
	if err := ch.Qos(PrefetchCount, 0, false); err != nil {
		log.Fatalf("[charging] set QoS: %v", err)
	}

	// Ensure topology exists (idempotent declarations).
	if err := ch.ExchangeDeclare(ExchangeCommands, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		log.Fatalf("[charging] declare commands exchange: %v", err)
	}
	if err := ch.ExchangeDeclare(ExchangeResponses, amqp.ExchangeDirect, true, false, false, false, nil); err != nil {
		log.Fatalf("[charging] declare responses exchange: %v", err)
	}
	if _, err := ch.QueueDeclare(QueueCharging, true, false, false, false, nil); err != nil {
		log.Fatalf("[charging] declare queue: %v", err)
	}
	if err := ch.QueueBind(QueueCharging, QueueCharging, ExchangeCommands, false, nil); err != nil {
		log.Fatalf("[charging] bind queue: %v", err)
	}
	if _, err := ch.QueueDeclare(QueueResponse, true, false, false, false, nil); err != nil {
		log.Fatalf("[charging] declare response queue: %v", err)
	}
	if err := ch.QueueBind(QueueResponse, QueueResponse, ExchangeResponses, false, nil); err != nil {
		log.Fatalf("[charging] bind response queue: %v", err)
	}

	deliveries, err := ch.Consume(
		QueueCharging,
		"charging-consumer",
		false, // autoAck = false
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[charging] consume: %v", err)
	}

	log.Println("[charging] service started, waiting for commands...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[charging] shutting down")
			return
		case d, ok := <-deliveries:
			if !ok {
				log.Println("[charging] channel closed")
				return
			}
			handleDelivery(ctx, ch, d)
		}
	}
}

func handleDelivery(ctx context.Context, ch *amqp.Channel, d amqp.Delivery) {
	var cmd SagaCommand
	if err := json.Unmarshal(d.Body, &cmd); err != nil {
		log.Printf("[charging] invalid message, discarding: %v", err)
		_ = d.Ack(false)
		return
	}

	log.Printf("[charging] received: saga=%s action=%s station=%s",
		cmd.SagaID, cmd.Action, cmd.StationID)

	var resp SagaResponse
	resp.SagaID = cmd.SagaID
	resp.Step = cmd.Step

	switch cmd.Action {
	case "EXECUTE":
		resp = startCharging(cmd)
	case "ROLLBACK":
		resp = stopCharging(cmd)
	default:
		log.Printf("[charging] unknown action: %s", cmd.Action)
		resp.Success = false
		resp.Error = fmt.Sprintf("unknown action: %s", cmd.Action)
	}

	if err := publishResponse(ctx, ch, resp); err != nil {
		log.Printf("[charging] failed to publish response: %v — nacking", err)
		_ = d.Nack(false, true)
		return
	}

	_ = d.Ack(false)
}

func startCharging(cmd SagaCommand) SagaResponse {
	resp := SagaResponse{SagaID: cmd.SagaID, Step: cmd.Step}

	// Simulate hardware initialization delay.
	time.Sleep(500 * time.Millisecond)

	// Simulate charging start — 95% success rate.
	if rand.Intn(100) < 95 {
		log.Printf("[charging] saga=%s: charging STARTED at station %s for user %s",
			cmd.SagaID, cmd.StationID, cmd.UserID)
		resp.Success = true
	} else {
		log.Printf("[charging] saga=%s: charging FAILED at station %s (hardware error, simulated)",
			cmd.SagaID, cmd.StationID)
		resp.Success = false
		resp.Error = "charger hardware error"
	}
	return resp
}

func stopCharging(cmd SagaCommand) SagaResponse {
	resp := SagaResponse{SagaID: cmd.SagaID, Step: cmd.Step}

	// Simulate stopping charge — always succeeds.
	log.Printf("[charging] saga=%s: charging STOPPED at station %s (compensating)",
		cmd.SagaID, cmd.StationID)
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
