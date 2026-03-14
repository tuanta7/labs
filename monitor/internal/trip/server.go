package trip

import (
	"context"
	"net/http"

	httpmiddleware "github.com/tuanta7/monitor/pkg/http"
	"github.com/tuanta7/monitor/pkg/monitor"
)

type Server struct {
	mux     *http.ServeMux
	server  *http.Server
	handler *Handler
	logger  *monitor.Logger
	meter   *monitor.Meter
	tracer  *monitor.Tracer
}

func NewServer(
	addr string,
	handler *Handler,
	logger *monitor.Logger,
	meter *monitor.Meter,
	tracer *monitor.Tracer,
) *Server {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &Server{
		mux:     mux,
		server:  server,
		handler: handler,
		logger:  logger,
		meter:   meter,
		tracer:  tracer,
	}
}

func (s *Server) Run() error {
	monitor.InitHTTPMeter(s.meter)

	s.mux.Handle("POST /trips", monitor.Middleware(s.tracer, s.logger,
		httpmiddleware.VerifyFakeToken(
			s.handler.CreateTrip,
		),
	))

	s.mux.Handle("POST /trips/{id}", httpmiddleware.VerifyFakeToken(s.handler.AcceptTrip))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
