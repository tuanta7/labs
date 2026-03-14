package main

import (
	"context"
	"net/http"

	"github.com/tuanta7/monitor/internal/trip"
	"github.com/tuanta7/monitor/pkg/graceful"
	"github.com/tuanta7/monitor/pkg/monitor"
	"github.com/tuanta7/monitor/pkg/slient"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	cfg, err := trip.LoadConfig()
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

	tracerProvider, err := monitor.NewTracerProvider(ctx, cfg.OTelServiceName, grpcConn)
	slient.PanicOnErr(err, "failed to create tracer provider")

	meterProvider, err := monitor.NewMeterProvider(ctx, cfg.OTelServiceName, grpcConn)
	slient.PanicOnErr(err, "failed to create meter provider")

	tracer := monitor.NewTracer(tracerProvider, cfg.OTelServiceName)
	meter := monitor.NewMeter(meterProvider, cfg.OTelServiceName)

	mongoClient, err := mongo.Connect(options.Client().
		ApplyURI(cfg.MongoConfig.URI).
		SetMaxPoolSize(cfg.MongoConfig.MaxPoolSize).
		SetMinPoolSize(cfg.MongoConfig.MinPoolSize).
		SetConnectTimeout(cfg.MongoConfig.ConnectTimeout),
	)
	slient.PanicOnErr(err)
	defer slient.CloseWithContext(mongoClient.Disconnect, ctx)

	httpClient := &http.Client{
		Timeout: cfg.HTTPClientTimeout,
	}

	repository := trip.NewRepository(mongoClient)
	business := trip.NewUseCase(repository, httpClient, logger, meter, tracer)
	handler := trip.NewHandler(business)
	server := trip.NewServer(cfg.BindAddress, handler, logger, meter, tracer)

	logger.Info("starting server", zap.String("address", cfg.BindAddress))
	err = graceful.RunServer(server)
	logger.Info("server stopped", zap.Error(err))
}
