package main

import (
	"context"
	"time"

	"github.com/tuanta7/k6noz/services/internal/ingestion"
	"github.com/tuanta7/k6noz/services/pkg/graceful"
	"github.com/tuanta7/k6noz/services/pkg/kafka"
	"github.com/tuanta7/k6noz/services/pkg/otelx"
	"github.com/tuanta7/k6noz/services/pkg/slient"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	cfg, err := ingestion.LoadConfig()
	slient.PanicOnErr(err)

	logger, err := zapx.NewLogger(zap.DebugLevel)
	slient.PanicOnErr(err)
	defer slient.Close(logger)

	prometheus, err := otelx.NewPrometheusProvider()
	slient.PanicOnErr(err)

	grpcConn, err := grpc.NewClient(cfg.OTelGRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	slient.PanicOnErr(err)

	monitor := otelx.NewMonitor(cfg.OTelServiceName, grpcConn, otelx.WithPrometheus(prometheus))
	defer slient.CloseWithContext(monitor, ctx)

	err = monitor.SetupOtelSDK(ctx)
	slient.PanicOnErr(err)

	publisher, err := kafka.NewPublisher(cfg.Kafka.Brokers)
	slient.PanicOnErr(err)
	defer publisher.Close()

	uc := ingestion.NewUseCase(logger, publisher)
	handler := ingestion.NewHandler(logger, uc)
	server := ingestion.NewServer(cfg, handler, prometheus.Handler())

	logger.Info("starting server", zap.String("address", cfg.BindAddress))
	err = graceful.RunServer(server, 20*time.Second)
	logger.Info("server stopped", zap.Error(err))
}
