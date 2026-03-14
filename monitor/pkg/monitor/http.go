package monitor

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
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

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (w *responseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func MetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		requestTotal.Add(r.Context(), 1, metric.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
			attribute.Int("status_code", rw.StatusCode),
		))

		if rw.StatusCode >= 400 {
			requestError.Add(r.Context(), 1, metric.WithAttributes(
				attribute.String("method", r.Method),
				attribute.String("path", r.URL.Path),
				attribute.Int("status_code", rw.StatusCode),
			))
		}
	})
}

func TraceMiddleware(tracer *Tracer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.StartNewSpan(r.Context(), "http.request")
		defer span.End()

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r.WithContext(ctx))

		span.SetAttributes(
			semconv.HTTPStatusCode(rw.StatusCode),
			semconv.HTTPMethod(r.Method),
			semconv.URLPath(r.URL.Path),
		)
	})
}

func Middleware(tracer *Tracer, logger *Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request received",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		handler := MetricMiddleware(TraceMiddleware(tracer, next))
		handler.ServeHTTP(w, r)
	})
}
