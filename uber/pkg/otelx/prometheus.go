package otelx

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
)

type PrometheusProvider struct {
	registry *prometheus.Registry
	exporter *otelprom.Exporter
}

func NewPrometheusProvider(cs ...prometheus.Collector) (*PrometheusProvider, error) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(cs...)

	exporter, err := otelprom.New(
		otelprom.WithRegisterer(registry),
	)
	if err != nil {
		return nil, err
	}

	return &PrometheusProvider{
		registry: registry,
		exporter: exporter,
	}, nil
}

func (p *PrometheusProvider) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		p.registry, promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}),
	)
}

func (p *PrometheusProvider) Exporter() *otelprom.Exporter {
	return p.exporter
}
