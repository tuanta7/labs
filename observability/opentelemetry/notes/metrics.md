# Metrics (OpenTelemetry)

Metrics are numerical time-series data collection that track the health and performance of a system over time. They help identify trends, anomalies, and overall system health. They can represent:

- **Infrastructure-level** Metrics (often collected by agents): CPU, Memory, Disk I/O, Network throughput, etc
- **Application-level** Metrics: Request counts per endpoint, Latency histograms, Cache hit/miss ratio, Queue length, etc
- **Custom** metrics: Number of orders placed, Items in shopping cart, Messages processed per second, etc.

## 1. Metrics API

The Metrics API consists of these main components:

- **MeterProvider**: The entry point of the API, holds meter configuration and provides access to Meters.
- **Meter**: Responsible for creating Instruments.
- **Instrument**: Responsible for reporting Measurements.
