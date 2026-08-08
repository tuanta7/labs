package o11y

import (
	"context"
	"net/http"
	"time"

	promcli "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelprom "go.opentelemetry.io/contrib/bridges/prometheus"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
)

type PrometheusMeterProvider struct {
	provider *metric.MeterProvider
	handler  http.Handler
}

func NewPrometheusMeterProvider(ctx context.Context, service string) (*PrometheusMeterProvider, error) {
	registry := promcli.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(grpcExportEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(resource.Default(),
		resource.NewSchemaless(semconv.ServiceName(service)))
	if err != nil {
		return nil, err
	}

	reader := metric.NewPeriodicReader(exporter,
		metric.WithProducer(otelprom.NewMetricProducer(otelprom.WithGatherer(registry))),
		metric.WithInterval(5*time.Second))

	return &PrometheusMeterProvider{
		provider: metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(reader),
		),
		handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	}, nil
}

func (p *PrometheusMeterProvider) Shutdown(ctx context.Context) error {
	if p.provider == nil {
		return nil
	}

	return p.provider.Shutdown(ctx)
}

func (p *PrometheusMeterProvider) MetricsHandler() http.Handler {
	return p.handler
}
