package main

import (
	"context"

	"github.com/tuanta7/monitor/internal/notification"
	"github.com/tuanta7/monitor/pkg/graceful"
	"github.com/tuanta7/monitor/pkg/monitor"
	"github.com/tuanta7/monitor/pkg/slient"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := notification.LoadConfig()
	slient.PanicOnErr(err)

	monitor.InitPropagator()

	logger, err := monitor.NewLogger()
	slient.PanicOnErr(err, "failed to create logger")
	defer slient.Close(logger)

	grpcConn, err := grpc.NewClient(cfg.OTelGRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	slient.PanicOnErr(err)
	defer slient.Close(grpcConn)

	tracerProvider, err := monitor.NewTracerProvider(context.Background(), cfg.OTelServiceName, grpcConn)
	slient.PanicOnErr(err, "failed to create tracer provider")

	uc := notification.NewUseCase(monitor.NewTracer(tracerProvider, cfg.OTelServiceName))
	handler := notification.NewHandler(uc)
	server := notification.NewServer(cfg, handler)

	logger.Info("starting server", zap.String("address", cfg.BindAddress))
	err = graceful.RunServer(server)
	logger.Info("server stopped", zap.Error(err))
}
