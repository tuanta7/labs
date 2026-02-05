package trip

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/kafka"
)

type Handler struct {
	kafka kafka.Consumer
	uc    *UseCase
}

func (h *Handler) InsertLocations(ctx context.Context) error {
	var location []*domain.Location
	err := h.uc.InsertLocations(ctx, location)
	if err != nil {
		return err
	}

	return nil
}
