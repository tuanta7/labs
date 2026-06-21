package trip

import (
	"net/http"
	"time"

	httpx "github.com/tuanta7/monitor/pkg/http"
)

type Handler struct {
	uc *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	profile, ok := r.Context().Value("profile").(*httpx.Profile)
	if !ok || profile == nil {
		_ = httpx.ErrorJSON(w, httpx.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "profile not found in context",
		})
		return
	}

	var pickUpLocation Location
	err := httpx.ReadJSON(r, &pickUpLocation)
	if err != nil {
		_ = httpx.ErrorJSON(w, httpx.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	if !pickUpLocation.Valid() {
		_ = httpx.ErrorJSON(w, httpx.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid pickup location value",
		})
		return
	}

	time.Sleep(10 * time.Millisecond)
	err = h.uc.CreateTrip(r.Context(), &Trip{
		PassengerID:    profile.ID,
		PickUpLocation: pickUpLocation,
	})
	if err != nil {
		_ = httpx.ErrorJSON(w, httpx.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to create trip",
		})
		return
	}
	time.Sleep(10 * time.Millisecond)

	_ = httpx.WriteJSON(w, http.StatusOK, httpx.JSON{})
}

func (h *Handler) AcceptTrip(w http.ResponseWriter, r *http.Request) {
	profile, ok := r.Context().Value("profile").(*httpx.Profile)
	if !ok || profile == nil {
		_ = httpx.ErrorJSON(w, httpx.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "profile not found in context",
		})
		return
	}

	tripID := r.PathValue("id")
	if tripID == "" {
		_ = httpx.ErrorJSON(w, httpx.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "trip id is required",
		})
		return
	}

	err := h.uc.AcceptTrip(r.Context(), tripID, profile.ID)
	if err != nil {
		_ = httpx.ErrorJSON(w, httpx.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to get trip",
		})
		return
	}

	_ = httpx.WriteJSON(w, http.StatusOK, httpx.JSON{})
}
