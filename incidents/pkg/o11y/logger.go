package o11y

import (
	"context"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	provider *sdklog.LoggerProvider
}

func NewLogger(ctx context.Context, service string) (*Logger, error) {
	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(grpcExportEndpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(resource.Default(),
		resource.NewSchemaless(semconv.ServiceName(service)))
	if err != nil {
		return nil, err
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)
	otelCore := otelzap.NewCore(service, otelzap.WithLoggerProvider(provider))

	localCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logger := zap.New(zapcore.NewTee(localCore, otelCore), zap.AddCaller())

	return &Logger{
		Logger:   logger,
		provider: provider,
	}, nil
}

func (l *Logger) Shutdown(ctx context.Context) error {
	if l.provider == nil {
		return nil
	}

	return l.provider.Shutdown(ctx)
}
