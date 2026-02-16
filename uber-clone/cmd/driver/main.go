package main

import (
	"context"
	"time"

	"github.com/tuanta7/k6noz/services/internal/driver"
	"github.com/tuanta7/k6noz/services/pkg/graceful"
	"github.com/tuanta7/k6noz/services/pkg/mongo"
	"github.com/tuanta7/k6noz/services/pkg/otelx"
	"github.com/tuanta7/k6noz/services/pkg/slient"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	cfg, err := driver.LoadConfig()
	slient.PanicOnErr(err)

	logger, err := zapx.NewLogger()
	slient.PanicOnErr(err, "failed to create logger")
	defer slient.Close(logger)

	grpcConn, err := grpc.NewClient(cfg.OTelGRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	slient.PanicOnErr(err)
	defer slient.Close(grpcConn)

	monitor := otelx.NewMonitor(cfg.OTelServiceName, grpcConn)
	defer slient.CloseWithContext(monitor, ctx)

	err = monitor.SetupOtelSDK(ctx)
	slient.PanicOnErr(err)

	mongoClient, err := mongo.NewClient(ctx, &mongo.Config{
		URI:            cfg.MongoConfig.URI,
		Database:       cfg.MongoConfig.Database,
		ConnectTimeout: cfg.MongoConfig.ConnectTimeout,
		QueryTimeout:   cfg.MongoConfig.QueryTimeout,
		Monitor:        true,
	})
	slient.PanicOnErr(err)
	defer slient.CloseWithContext(mongoClient, ctx)

	repo := driver.NewRepository(mongoClient)
	uc := driver.NewUseCase(logger, repo)
	handler := driver.NewHandler(logger, uc)

	server := driver.NewServer(cfg.BindAddress, handler)

	logger.Info("starting server", zap.String("address", cfg.BindAddress))
	err = graceful.RunServer(server, 20*time.Second)
	logger.Info("server stopped", zap.Error(err))
}
