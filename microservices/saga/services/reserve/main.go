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
	QueueReservation  = "reservation"
	QueueResponse     = "response"
	PrefetchCount     = 1
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rabbitURL := envOrDefault("RABBITMQ_URL", "amqp://rabbitmq:password@localhost:5672/")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("[reserve] connect to rabbitmq: %v", err)
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[reserve] open channel: %v", err)
	}
	defer func() { _ = ch.Close() }()

	// Set QoS — process one message at a time.
	if err := ch.Qos(PrefetchCount, 0, false); err != nil {
		log.Fatalf("[reserve] set QoS: %v", err)
	}

	// Ensure topology exists (idempotent declarations).
	if err := ch.ExchangeDeclare(ExchangeCommands, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		log.Fatalf("[reserve] declare commands exchange: %v", err)
	}
	if err := ch.ExchangeDeclare(ExchangeResponses, amqp.ExchangeDirect, true, false, false, false, nil); err != nil {
		log.Fatalf("[reserve] declare responses exchange: %v", err)
	}
	if _, err := ch.QueueDeclare(QueueReservation, true, false, false, false, nil); err != nil {
		log.Fatalf("[reserve] declare queue: %v", err)
	}
	if err := ch.QueueBind(QueueReservation, QueueReservation, ExchangeCommands, false, nil); err != nil {
		log.Fatalf("[reserve] bind queue: %v", err)
	}
	if _, err := ch.QueueDeclare(QueueResponse, true, false, false, false, nil); err != nil {
		log.Fatalf("[reserve] declare response queue: %v", err)
	}
	if err := ch.QueueBind(QueueResponse, QueueResponse, ExchangeResponses, false, nil); err != nil {
		log.Fatalf("[reserve] bind response queue: %v", err)
	}

	deliveries, err := ch.Consume(
		QueueReservation,
		"reserve-consumer",
		false, // autoAck = false
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[reserve] consume: %v", err)
	}

	log.Println("[reserve] service started, waiting for commands...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[reserve] shutting down")
			return
		case d, ok := <-deliveries:
			if !ok {
				log.Println("[reserve] channel closed")
				return
			}
			handleDelivery(ctx, ch, d)
		}
	}
}

func handleDelivery(ctx context.Context, ch *amqp.Channel, d amqp.Delivery) {
	var cmd SagaCommand
	if err := json.Unmarshal(d.Body, &cmd); err != nil {
		log.Printf("[reserve] invalid message, discarding: %v", err)
		_ = d.Ack(false)
		return
	}

	log.Printf("[reserve] received: saga=%s action=%s station=%s", cmd.SagaID, cmd.Action, cmd.StationID)

	var resp SagaResponse
	resp.SagaID = cmd.SagaID
	resp.Step = cmd.Step

	switch cmd.Action {
	case "EXECUTE":
		resp = executeReservation(cmd)
	case "ROLLBACK":
		resp = rollbackReservation(cmd)
	default:
		log.Printf("[reserve] unknown action: %s", cmd.Action)
		resp.Success = false
		resp.Error = fmt.Sprintf("unknown action: %s", cmd.Action)
	}

	if err := publishResponse(ctx, ch, resp); err != nil {
		log.Printf("[reserve] failed to publish response: %v — nacking", err)
		_ = d.Nack(false, true)
		return
	}

	_ = d.Ack(false)
}

func executeReservation(cmd SagaCommand) SagaResponse {
	resp := SagaResponse{
		SagaID: cmd.SagaID,
		Step:   cmd.Step,
	}

	// Simulate reservation logic.
	if rand.Intn(10) < 9 {
		// 90% success rate for demonstration purposes.
		log.Printf("[reserve] saga=%s: station %s reserved successfully", cmd.SagaID, cmd.StationID)
		resp.Success = true
	} else {
		log.Printf("[reserve] saga=%s: station %s reservation FAILED (simulated)", cmd.SagaID, cmd.StationID)
		resp.Success = false
		resp.Error = "station unavailable"
	}
	return resp
}

func rollbackReservation(cmd SagaCommand) SagaResponse {
	resp := SagaResponse{
		SagaID: cmd.SagaID,
		Step:   cmd.Step,
	}
	// Simulate cancellation — always succeeds.
	log.Printf("[reserve] saga=%s: station %s reservation CANCELLED (compensating)", cmd.SagaID, cmd.StationID)
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
