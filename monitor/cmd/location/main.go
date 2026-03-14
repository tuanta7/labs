package main

import (
	"github.com/tuanta7/monitor/internal/location"
	"github.com/tuanta7/monitor/pkg/monitor"
	"github.com/tuanta7/monitor/pkg/slient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//ctx := context.Background()

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
}
