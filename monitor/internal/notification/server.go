package notification

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Server struct {
	server  *http.Server
	mux     *echo.Echo
	handler *Handler
}

func NewServer(cfg *Config, handler *Handler) *Server {
	mux := echo.New()
	return &Server{
		handler: handler,
		mux:     mux,
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
	s.mux.GET("/metrics", nil)
	s.mux.POST("/push", s.handler.SendPushNotification)
}
