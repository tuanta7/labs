package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Meter struct {
	service string
	metric.Meter
}

func NewMeter(provider *sdkmetric.MeterProvider, service string) *Meter {
	meter := provider.Meter(service)
	return &Meter{
		service: service,
		Meter:   meter,
	}
}