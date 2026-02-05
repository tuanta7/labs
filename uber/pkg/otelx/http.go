package otelx

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Handler(handler func(w http.ResponseWriter, r *http.Request), op string) http.Handler {
	return otelhttp.NewHandler(http.HandlerFunc(handler), op)
}
