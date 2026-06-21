package monitor

import (
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
