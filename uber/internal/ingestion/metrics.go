package ingestion

import (
	"github.com/tuanta7/k6noz/services/pkg/otelx"
	"go.opentelemetry.io/otel"
)

const (
	ServiceName = "services/internal/ingestion"
)

var (
	meter = otel.Meter(ServiceName)

	locationUpdatesValidTotal   = otelx.MustMetric(meter.Int64Counter("location_updates_valid_total"))
	locationUpdatesInvalidTotal = otelx.MustMetric(meter.Int64Counter("location_updates_invalid_total"))
)
