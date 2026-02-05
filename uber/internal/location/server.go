package location

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	address string
	engine  *echo.Echo
	server  *http.Server
	handler *Handler
}

func NewServer(address string, handler *Handler) *Server {
	engine := echo.New()
	server := &http.Server{
		Addr:    address,
		Handler: engine,
	}

	return &Server{
		address: address,
		engine:  engine,
		server:  server,
		handler: handler,
	}
}

func (s *Server) Run() error {
	s.registerRoutes()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerRoutes() {}
