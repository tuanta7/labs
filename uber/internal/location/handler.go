package location

import (
	"github.com/labstack/echo/v4"
	"github.com/tuanta7/k6noz/services/pkg/kafka"
)

type Handler struct {
	kafka kafka.Consumer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetNearByDrivers(ctx echo.Context) error {
	return nil
}

func (h *Handler) GetDriverLocation() error {
	// sse
	return nil
}

func (h *Handler) UpdateDriverLocation() error {
	// kafka
	return nil
}
