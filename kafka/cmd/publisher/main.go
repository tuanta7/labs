package main

import (
	"kafka-lab/internal/config"
	"kafka-lab/internal/handler"
	"kafka-lab/internal/kafka"
	"kafka-lab/internal/usecase"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	zl "github.com/rs/zerolog/log"
)

func main() {
	logger := zl.Logger.With().Timestamp().Logger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
	}

	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create kafka producer")
	}
	defer producer.Close()

	logger.Info().
		Strs("brokers", cfg.KafkaBrokers).
		Str("topic", cfg.KafkaTopicLocation).
		Msg("kafka producer initialized")

	uc := usecase.NewLocationUC(&logger, usecase.WithPublisher(producer, cfg.KafkaTopicLocation))
	hdl := handler.NewPublishHandler(uc, &logger)

	app := fiber.New()
	app.Get("/ws", websocket.New(hdl.Handle))

	err = app.Listen(":3000")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}
}
