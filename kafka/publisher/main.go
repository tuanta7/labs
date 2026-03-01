package main

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	zl "github.com/rs/zerolog/log"
)

func main() {
	logger := zl.Logger.With().Timestamp().Logger()

	cfg, err := LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
	}

	producer, err := NewProducer(cfg.KafkaBrokers)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create kafka producer")
	}
	defer producer.Close()

	logger.Info().
		Strs("brokers", cfg.KafkaBrokers).
		Str("topic", cfg.KafkaTopicLocation).
		Msg("kafka producer initialized")

	uc := NewUseCase(producer, cfg.KafkaTopicLocation, &logger)
	handler := NewPublishHandler(uc, &logger)

	app := fiber.New()
	app.Get("/ws", websocket.New(handler.Handle))

	err = app.Listen(":3000")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}
}
