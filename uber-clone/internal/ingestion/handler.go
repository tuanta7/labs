package ingestion

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer h.closeWS(ws)

	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		msgType, msg, err := ws.ReadMessage()
		h.logger.Debug("received message", zap.Int("msgType", msgType))
		if err != nil {
			h.logger.Debug("websocket read error", zap.Error(err))
			break
		}

		var location domain.DriverLocationMessage
		err = json.Unmarshal(msg, &location)
		if err != nil {
			h.logger.Debug("invalid payload", zap.Error(err))
			locationUpdatesInvalidTotal.Add(r.Context(), 1)
			continue
		}

		locationUpdatesValidTotal.Add(r.Context(), 1)
		h.uc.PublishLocation(r.Context(), &location)
	}
}

func (h *Handler) closeWS(ws *websocket.Conn) {
	if closeErr := ws.Close(); closeErr != nil {
		h.logger.Warn("websocket close error", zap.Error(closeErr))
	}
}
