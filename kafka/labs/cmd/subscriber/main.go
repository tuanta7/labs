package main

import (
	"context"
	"kafka-lab/internal/config"
	"kafka-lab/internal/handler"
	"kafka-lab/internal/kafka"
	"sync"

	"kafka-lab/internal/repository"
	"kafka-lab/internal/usecase"

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

	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopicLocation, cfg.KafkaConsumerGroup)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create kafka consumer")
	}
	defer consumer.Close()

	repo := repository.NewRepository(session)
	uc := usecase.NewLocationUC(&logger, usecase.WithRepository(repo))
	hdl := handler.NewConsumeHandler(consumer, uc, &logger)

	var handlers []func(ctx context.Context) error
	handlers = append(handlers, hdl.ConsumeLocationMessage)

	var wg sync.WaitGroup
	for _, h := range handlers {
		wg.Go(func() {
			consumeErr := h(context.Background())
			if consumeErr != nil {
				logger.Fatal().Err(consumeErr).Msg("failed to consume message")
			}
		})
	}

	wg.Wait()
}
