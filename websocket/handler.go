package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
	writeWait       = 10 * time.Second
	pongWait        = 10 * time.Second
	pingPeriod      = (pongWait * 9) / 10
)

type Handler struct {
	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) HandlePolling(w http.ResponseWriter, r *http.Request) {}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBufferSize,
	WriteBufferSize: writeBufferSize,
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
	defer ws.Close()

	h.logger.Info("websocket connection established",
		zap.String("remoteAddr", ws.RemoteAddr().String()),
	)

	done := make(chan struct{})

	go func() {
		defer close(done)

		_ = ws.SetReadDeadline(time.Now().Add(pongWait))
		ws.SetPongHandler(func(string) error {
			h.logger.Debug("pong received from client")
			_ = ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		for {
			msgType, msg, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					h.logger.Error("websocket unexpected close", zap.Error(err))
				} else {
					h.logger.Debug("websocket read error", zap.Error(err))
				}
				break
			}

			h.logger.Info("websocket message received",
				zap.Int("type", msgType),
				zap.ByteString("message", msg),
			)

			if err := ws.WriteMessage(msgType, msg); err != nil {
				h.logger.Debug("websocket write error", zap.Error(err))
				break
			}
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send protocol-level ping frame
			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				h.logger.Debug("failed to send ping frame", zap.Error(err))
				return
			}
			h.logger.Debug("protocol ping frame sent to client")

			// Also send application-level PING text message
			if err := ws.WriteMessage(websocket.TextMessage, []byte("PING")); err != nil {
				h.logger.Debug("failed to send PING text message", zap.Error(err))
				return
			}
			h.logger.Info("PING text message sent to client")
		case <-done:
			h.logger.Info("websocket connection closed",
				zap.String("remoteAddr", ws.RemoteAddr().String()),
			)
			return
		}
	}
}

func (h *Handler) HandleBroadcast(w http.ResponseWriter, r *http.Request) {}

// Connect send a connection request to the server
func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ServerURL string `json:"serverUrl"`
	}

	err := ReadJSON(r, &requestPayload)
	if err != nil {
		_ = ErrorJSON(w, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
}
