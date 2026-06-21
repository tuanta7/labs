package notification

import (
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Handler struct {
	uc *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) SendPushNotification(c *echo.Context) error {
	ctx := otel.GetTextMapPropagator().Extract(
		c.Request().Context(),
		propagation.HeaderCarrier(c.Request().Header),
	)
	return h.uc.SendPushNotification(ctx)
}
