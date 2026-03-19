package handler

import (
	"context"
	"kafka-lab/internal/usecase"

	"github.com/rs/zerolog"
)

type SubscribeHandler struct {
	uc     *usecase.ConsumerUC
	logger *zerolog.Logger
}

func NewSubscribeHandler(uc *usecase.ConsumerUC, logger *zerolog.Logger) *SubscribeHandler {
	if logger == nil {
		logger = zerolog.DefaultContextLogger
	}

	return &SubscribeHandler{
		uc:     uc,
		logger: logger,
	}
}

func (h *SubscribeHandler) ConsumeLocationTopic(ctx context.Context) error {
	if err := h.uc.ConsumeLocationTopic(ctx); err != nil {
		h.logger.Error().
			Err(err).
			Msg("failed to consume location topic")
		return err
	}

	return nil
}
