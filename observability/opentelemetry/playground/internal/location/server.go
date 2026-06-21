package location

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tuanta7/monitor/pkg/monitor"
	"go.uber.org/zap"
)

type Server struct {
	*http.Server
	echo   *echo.Echo
	logger *zap.Logger
	meter  *monitor.Meter
	tracer *monitor.Tracer
}

func NewServer(
	echo *echo.Echo,
	addr string,
	logger *zap.Logger,
	meter *monitor.Meter,
	tracer *monitor.Tracer,
) *Server {
	s := &Server{
		echo:   echo,
		logger: logger,
		meter:  meter,
		tracer: tracer,
	}
	s.Server = &http.Server{
		Addr:    addr,
		Handler: echo,
	}
	return s
}
