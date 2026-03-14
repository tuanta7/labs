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

type PrometheusMeter struct {
	registry     *prometheus.Registry
	exporter     *otelprom.Exporter
	HttpRequests metric.Int64Counter
}

func NewPrometheusMeter(cs ...prometheus.Collector) (*PrometheusMeter, error) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(cs...)

	exporter, err := otelprom.New(
		otelprom.WithRegisterer(registry),
	)
	if err != nil {
		return nil, err
	}

	meter := otel.GetMeterProvider().Meter("http-server")
	httpRequests, err := meter.Int64Counter("http.server.requests", metric.WithDescription("Number of HTTP requests"))
	if err != nil {
		return nil, err
	}

	return &PrometheusMeter{
		registry:     registry,
		exporter:     exporter,
		HttpRequests: httpRequests,
	}, nil
}

func (m *PrometheusMeter) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		m.registry, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
	)
}
