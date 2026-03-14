package monitor

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	requestTotal metric.Int64Counter
	requestError metric.Int64Counter
)

func InitHTTPMeter(meter *Meter) {
	requestTotal, _ = meter.Int64Counter("http.server.requests",
		metric.WithDescription("Number of HTTP requests"),
	)

	requestError, _ = meter.Int64Counter("http.server.errors",
		metric.WithDescription("Number of HTTP errors"),
	)
}

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func MetricMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTotal.Add(r.Context(), 1)
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		if rw.StatusCode >= 400 {
			requestError.Add(r.Context(), 1)
		}
	}
}

func TraceMiddleware(tracer *Tracer, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.StartNewSpan(r.Context(), "http.request")

		span.SetAttributes(attribute.KeyValue{
			Key:   "http_method",
			Value: attribute.StringValue(r.Method),
		}, attribute.KeyValue{
			Key:   "path",
			Value: attribute.StringValue(r.URL.Path),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
		span.End()
	}
}

func Middleware(tracer *Tracer, logger *Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request received",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		MetricMiddleware(TraceMiddleware(tracer, next))
	}
}
