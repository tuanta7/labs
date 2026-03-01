package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/rs/zerolog"
)

type PublishHandler struct {
	uc     *LocationUseCase
	logger *zerolog.Logger
}

func NewPublishHandler(uc *LocationUseCase, logger *zerolog.Logger) *PublishHandler {
	return &PublishHandler{
		uc:     uc,
		logger: logger,
	}
}

func (h *PublishHandler) Handle(c *websocket.Conn) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to read websocket message")
			break
		}

		var location Location
		if err = json.Unmarshal(msg, &location); err != nil {
			h.logger.Error().Err(err).Msg("failed to unmarshal location data")
			continue
		}

		h.logger.Info().
			Str("location_id", location.ID).
			Str("user_id", location.UserID).
			Msg("received location from websocket")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		h.uc.ProduceLocation(ctx, &location)
		cancel()
	}
}
