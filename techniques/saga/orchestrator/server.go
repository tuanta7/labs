package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	ch     *amqp.Channel
	store  *SagaStore
	mux    *http.ServeMux
	server *http.Server
}

func NewServer(ch *amqp.Channel, store *SagaStore) *Server {
	return &Server{
		mux:   http.NewServeMux(),
		ch:    ch,
		store: store,
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	log.Println("[orchestrator] shutting down...")
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Println("[orchestrator] server shutdown: ", err)
	}
}

func (s *Server) Start() {
	s.mux.HandleFunc("POST /saga/start", func(w http.ResponseWriter, r *http.Request) {
		var req ChargingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if req.StationID == "" || req.UserID == "" || req.Amount <= 0 {
			http.Error(w, "user_id, station_id, and a positive amount are required", http.StatusBadRequest)
			return
		}

		saga, err := StartSaga(r.Context(), s.ch, s.store, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to start saga: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(saga)
	})

	s.mux.HandleFunc("GET /saga/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		saga, err := s.store.GetSaga(r.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf("saga not found: %v", err), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(saga)
	})

	addr := getEnv("HTTP_ADDR", ":8080")
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	log.Printf("[orchestrator] HTTP server listening on %s", addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http server: %v", err)
	}
}
