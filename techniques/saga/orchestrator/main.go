package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	defaultMongoURI    = "mongodb://localhost:27017"
	defaultRabbitMQURL = "amqp://rabbitmq:password@localhost:5672/"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mongoURI := getEnv("MONGO_URI", defaultMongoURI)
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("connect to mongodb: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("disconnect mongodb: %v", err)
		}
	}()

	store := NewSagaStore(mongoClient)

	conn, err := amqp.Dial(getEnv("RABBITMQ_URL", defaultRabbitMQURL))
	if err != nil {
		log.Fatalf("connect to rabbitmq: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("open rabbitmq channel: %v", err)
	}
	defer ch.Close()

	if err = SetupTopology(ch); err != nil {
		log.Fatalf("setup topology: %v", err)
	}
	log.Println("[orchestrator] RabbitMQ topology ready")

	go func() {
		if err := ConsumeResponses(ctx, ch, store); err != nil {
			log.Printf("[orchestrator] consumer stopped: %v", err)
			cancel()
		}
	}()

	srv := NewServer(ch, store)
	go srv.Start()

	<-ctx.Done()
	srv.Shutdown(context.Background())
}
