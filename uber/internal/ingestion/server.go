package ingestion

import (
	"context"
	"net/http"
)

type Server struct {
	server        *http.Server
	mux           *http.ServeMux
	handler       *Handler
	metricHandler http.Handler
}

func NewServer(cfg *Config, handler *Handler, metricHandler http.Handler) *Server {
	mux := http.NewServeMux()

	return &Server{
		handler:       handler,
		metricHandler: metricHandler,
		mux:           mux,
		server: &http.Server{
			Addr:    cfg.BindAddress,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	s.registerRoutes()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	s.mux.Handle("GET /metrics", s.metricHandler)
	s.mux.HandleFunc("/ws", s.handler.HandleWS)
}
