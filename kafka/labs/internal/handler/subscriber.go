package handler

import (
	"context"
	"kafka-lab/internal/kafka"
	"kafka-lab/internal/usecase"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

type ConsumeHandler struct {
	consumer Consumer
	uc       *usecase.LocationUC
	logger   *zerolog.Logger
}

func NewConsumeHandler(consumer Consumer, uc *usecase.LocationUC, logger *zerolog.Logger) *ConsumeHandler {
	return &ConsumeHandler{
		consumer: consumer,
		uc:       uc,
		logger:   logger,
	}
}

func (c *ConsumeHandler) ConsumeLocationMessage(ctx context.Context) error {
	for {
		err := c.consumer.Consume(ctx, func(ctx context.Context, msg *kgo.Record) error {
			c.logger.Info().
				Str("topic", msg.Topic).
				Bytes("key", msg.Key).
				Int64("timestamp", msg.Timestamp.Unix()).
				Msg("received message from kafka")

			return c.uc.SaveLocation(ctx, msg.Value)
		})

		if err != nil {
			return err
		}
	}
}
