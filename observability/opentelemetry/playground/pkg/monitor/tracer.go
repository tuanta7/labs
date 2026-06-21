package monitor

import (
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	service string
	tracer  trace.Tracer
}

func NewTracer(provider *sdktrace.TracerProvider, service string) *Tracer {
	return &Tracer{
		service: service,
		tracer:  provider.Tracer(service),
	}
}

func (t *Tracer) Service() string {
	return t.service
}

func (t *Tracer) StartNewSpan(
	ctx context.Context,
	name string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}
