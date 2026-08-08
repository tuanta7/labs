package o11y

import "os"

var (
	grpcExportEndpoint = readEnv("OTEL_GRPC_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
)

func readEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
