package main

import (
	"context"
	"errors"
	"kafka-lab/internal/config"
	"kafka-lab/internal/handler"
	"kafka-lab/internal/kafka"
	"kafka-lab/internal/repository"
	"kafka-lab/internal/usecase"
	"os/signal"
	"syscall"

	"github.com/gocql/gocql"
	zl "github.com/rs/zerolog/log"
)

func main() {
	logger := zl.Logger.With().Timestamp().Logger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
	}

	cluster := gocql.NewCluster(cfg.ScyllaHosts...)
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create scylla session")
	}
	defer session.Close()

	subscriber, err := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopicLocation, cfg.KafkaConsumerGroup)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create kafka consumer")
	}
	defer subscriber.Close()

	repo := repository.NewRepository(session)
	uc := usecase.NewConsumerUseCase(subscriber, repo, cfg.KafkaTopicLocation, cfg.KafkaConsumerGroup, &logger)
	h := handler.NewSubscribeHandler(uc, &logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info().
		Strs("brokers", cfg.KafkaBrokers).
		Str("topic", cfg.KafkaTopicLocation).
		Str("group_id", cfg.KafkaConsumerGroup).
		Msg("starting kafka subscriber")

	err = h.ConsumeLocationTopic(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Fatal().Err(err).Msg("subscriber stopped with error")
	}
}
