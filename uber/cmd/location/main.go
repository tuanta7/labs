package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/tuanta7/k6noz/services/internal/location"
	"github.com/tuanta7/k6noz/services/pkg/redis"
	"github.com/tuanta7/k6noz/services/pkg/slient"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger, err := zapx.NewLogger()
	panicOnErr(err)
	defer slient.Close(logger)

	redisClient, err := redis.NewFailoverClient(ctx, &redis.Config{},
		redis.WithMetrics(),
		redis.WithTraces(),
	)
	panicOnErr(err)
	defer slient.Close(redisClient)

	handler := location.NewHandler()

	echoEngine := echo.New()
	echoEngine.GET("/nearby", handler.GetNearByDrivers)
	defer slient.Close(echoEngine)

	err = echoEngine.Start(":8080")
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
