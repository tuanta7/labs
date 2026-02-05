package driver

import (
	"context"
	"net/http"

	"github.com/tuanta7/k6noz/services/pkg/otelx"
)

type Server struct {
	mux     *http.ServeMux
	server  *http.Server
	handler *Handler
}

func NewServer(addr string, handler *Handler) *Server {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &Server{
		mux:     mux,
		server:  server,
		handler: handler,
	}
}

func (s *Server) Run() error {
	s.mux.Handle("GET /drivers/{id}", otelx.Handler(s.handler.GetDriverByID, "GetDriverByID"))
	s.mux.Handle("POST /ratings", otelx.Handler(s.handler.CreateNewRating, "CreateNewRating"))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
