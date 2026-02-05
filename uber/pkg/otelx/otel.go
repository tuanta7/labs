package otelx

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
	"google.golang.org/grpc"
)

type ShutdownFn func(context.Context) error

type Monitor struct {
	serviceName   string
	grpcConn      *grpc.ClientConn
	prometheus    *PrometheusProvider
	shutdownFuncs []ShutdownFn
}

type Option func(*Monitor)

func WithPrometheus(prometheus *PrometheusProvider) Option {
	return func(monitor *Monitor) {
		monitor.prometheus = prometheus
	}
}

func NewMonitor(serviceName string, grpcConn *grpc.ClientConn, opts ...Option) *Monitor {
	m := &Monitor{
		serviceName: serviceName,
		grpcConn:    grpcConn,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Monitor) SetupOtelSDK(ctx context.Context) error {
	msf, err := m.initMeterProvider(ctx)
	if err != nil {
		return err
	}

	tsf, err := m.initTracerProvider(ctx)
	if err != nil {
		return err
	}

	m.shutdownFuncs = []ShutdownFn{msf, tsf}
	return nil
}

func (m *Monitor) Close(ctx context.Context) (errs error) {
	for _, sf := range m.shutdownFuncs {
		if err := sf(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

func (m *Monitor) initMeterProvider(ctx context.Context) (ShutdownFn, error) {
	res, err := resource.New(ctx,
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.serviceName),
		))
	if err != nil {
		return nil, err
	}

	otlpExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithGRPCConn(m.grpcConn),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	var reader sdkmetric.Reader
	if m.prometheus != nil {
		reader = m.prometheus.Exporter()
	} else {
		reader = sdkmetric.NewPeriodicReader(otlpExporter)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	otel.SetMeterProvider(meterProvider)
	return meterProvider.Shutdown, nil
}

func (m *Monitor) initTracerProvider(ctx context.Context) (ShutdownFn, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.serviceName),
		))
	if err != nil {
		return nil, err
	}

	otlpExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(m.grpcConn),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpExporter, sdktrace.WithBatchTimeout(5*time.Second)),
	)

	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}
