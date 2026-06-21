# Traces (OpenTelemetry)

Traces represent the flow of a request as it moves through different components or services of a system. They show the path of a request, the time spent in each service, and any errors encountered along the way, making it easier to pinpoint bottlenecks and root causes of issues.

## 1. Tracing API

Reference: [Tracing API](https://opentelemetry.io/docs/specs/otel/trace/api/)

The Tracing API consists of these main components:

- **TracerProvider**: The entry point of the API. It provides access to Tracers.
- **Tracer**: Responsible for creating Spans.
- **Span**: The API to trace an operation.

## 2. Propagators

Reference: [Propagators API](https://opentelemetry.io/docs/specs/otel/context/api-propagators/)

When a request moves from one service to another, the trace context needs to be carried along. Propagators enable the continuation of a trace across different services or processes in a distributed system.
