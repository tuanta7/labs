package driver

import (
	"net/http"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zapx.Logger
	uc     *UseCase
}

func NewHandler(logger *zapx.Logger, uc *UseCase) *Handler {
	return &Handler{
		logger: logger,
		uc:     uc,
	}
}

func (h *Handler) GetDriverByID(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("id")
	if driverID == "" {
		h.logger.Error("driver id is required", zap.String("url", r.URL.String()))
		_ = ErrorJSON(w, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "driver id is required",
		})
		return
	}

	driver, err := h.uc.GetDriverByID(r.Context(), driverID)
	if err != nil {
		h.logger.Error("failed to get driver", zap.Error(err))
		_ = ErrorJSON(w, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to get driver",
		})
		return
	}

	_ = WriteJSON(w, http.StatusOK, driver)
}

func (h *Handler) CreateNewRating(w http.ResponseWriter, r *http.Request) {
	rating := &domain.Rating{}
	if err := ReadJSON(r, rating); err != nil {
		h.logger.Error("failed to read rating input", zap.Error(err))
		_ = ErrorJSON(w, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to read rating",
			Details: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}
}
