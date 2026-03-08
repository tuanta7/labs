package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"orchestrator/internal/config"
	"orchestrator/internal/saga"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	ch     *amqp.Channel
	store  *saga.Store
	mux    *http.ServeMux
	server *http.Server
}

func NewServer(ch *amqp.Channel, store *saga.Store) *Server {
	return &Server{
		mux:   http.NewServeMux(),
		ch:    ch,
		store: store,
	}
}

func (sv *Server) Shutdown(ctx context.Context) {
	log.Println("[orchestrator] shutting down...")
	err := sv.server.Shutdown(ctx)
	if err != nil {
		log.Println("[orchestrator] server shutdown: ", err)
	}
}

func (sv *Server) Start() {
	sv.mux.HandleFunc("POST /saga/charging", func(w http.ResponseWriter, r *http.Request) {
		var req ChargingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		if req.StationID == "" || req.UserID == "" || req.Amount <= 0 {
			http.Error(w, "user_id, station_id, and a positive amount are required", http.StatusBadRequest)
			return
		}

		sg, err := sv.StartChargingSaga(r.Context(), req)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to start saga: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(sg)
	})

	sv.mux.HandleFunc("GET /saga/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		sg, err := sv.store.GetSaga(r.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf("saga not found: %v", err), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sg)
	})

	addr := config.GetEnv("HTTP_ADDR", ":8080")
	sv.server = &http.Server{
		Addr:    addr,
		Handler: sv.mux,
	}

	log.Printf("[orchestrator] HTTP server listening on %s", addr)
	if err := sv.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http server: %v", err)
	}
}
