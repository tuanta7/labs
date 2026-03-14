package notification

import (
	"math/rand"
	"time"

	"github.com/labstack/echo/v5"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SendPushNotification(c *echo.Context) error {
	r := rand.Intn(10)
	if r > 9 {
		return echo.NewHTTPError(500, "failed to send push notification")
	}

	time.Sleep(2 * time.Second)
	return nil
}
