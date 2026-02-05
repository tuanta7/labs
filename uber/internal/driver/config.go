package driver

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/gommon/log"
)

const ConfigPrefix = "DRIVER"

type Config struct {
	BindAddress      string      `envconfig:"BIND_ADDRESS" required:"true"`
	OTelGRPCEndpoint string      `envconfig:"OTEL_GRPC_ENDPOINT" required:"true"`
	OTelServiceName  string      `envconfig:"OTEL_SERVICE_NAME" required:"true" default:"driver-service"`
	MongoConfig      MongoConfig `envconfig:"MONGO"`
}

type MongoConfig struct {
	URI            string        `envconfig:"URI" required:"true" default:"mongodb://localhost:27017"`
	Database       string        `envconfig:"DATABASE" required:"true"`
	MaxPoolSize    uint64        `envconfig:"MAX_POOL_SIZE" default:"100"`
	MinPoolSize    uint64        `envconfig:"MIN_POOL_SIZE" default:"10"`
	ConnectTimeout time.Duration `envconfig:"CONNECT_TIMEOUT" default:"10s"`
	QueryTimeout   time.Duration `envconfig:"QUERY_TIMEOUT" default:"10s"`
}

func LoadConfig(envFiles ...string) (*Config, error) {
	var cfg Config
	if err := godotenv.Load(envFiles...); err != nil {
		log.Warnf("failed to load .env file: %v", err)
	}

	if err := envconfig.Process(ConfigPrefix, &cfg); err != nil {
		return nil, fmt.Errorf("failed to process env variables: %w", err)
	}

	return &cfg, nil
}
