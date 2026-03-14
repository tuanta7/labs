package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tuanta7/monitor/internal/location"
	"github.com/tuanta7/monitor/pkg/graceful"
	"github.com/tuanta7/monitor/pkg/monitor"
	"github.com/tuanta7/monitor/pkg/slient"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	cfg, err := location.LoadConfig()
	slient.PanicOnErr(err)

	logger, err := monitor.NewLogger()
	slient.PanicOnErr(err, "failed to create logger")
	defer slient.Close(logger)

	grpcConn, err := grpc.NewClient(cfg.OTelGRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	slient.PanicOnErr(err)
	defer slient.Close(grpcConn)

	tracerProvider, err := monitor.NewTracerProvider(ctx, cfg.OTelServiceName, grpcConn)
	slient.PanicOnErr(err, "failed to create tracer provider")

	meterProvider, err := monitor.NewMeterProvider(ctx, cfg.OTelServiceName, grpcConn)
	slient.PanicOnErr(err, "failed to create meter provider")

	tracer := monitor.NewTracer(tracerProvider, cfg.OTelServiceName)
	meter := monitor.NewMeter(meterProvider, cfg.OTelServiceName)

	echoEngine := echo.New()
	server := location.NewServer(echoEngine, cfg.BindAddress, logger, meter, tracer)

	logger.Info("starting server", zap.String("address", cfg.BindAddress))
	err = graceful.RunServer(server)
	logger.Info("server stopped", zap.Error(err))
}
