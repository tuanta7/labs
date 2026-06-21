package monitor

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

type Prometheus struct {
	registry *prometheus.Registry
}

func NewPrometheus(cs ...prometheus.Collector) (*Prometheus, error) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(cs...)

	return &Prometheus{
		registry: registry,
	}, nil
}

func (p *Prometheus) Registry() *prometheus.Registry {
	return p.registry
}

func (p *Prometheus) NewMeterProvider(ctx context.Context, service string) (*sdkmetric.MeterProvider, error) {
	res, err := resource.New(ctx,
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(service),
		))
	if err != nil {
		return nil, err
	}

	exporter, err := otelprom.New(
		otelprom.WithRegisterer(p.registry),
	)
	if err != nil {
		return nil, err
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	), nil
}

func (p *Prometheus) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		p.registry, promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}),
	)
}
