package location

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BindAddress      string `envconfig:"BIND_ADDRESS" default:":8081"`
	OTelGRPCEndpoint string `envconfig:"OTEL_GRPC_ENDPOINT" default:"localhost:4317"`
	OTelServiceName  string `envconfig:"OTEL_SERVICE_NAME" default:"location"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
