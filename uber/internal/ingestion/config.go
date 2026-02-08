package ingestion

import (
	"fmt"

	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const ConfigPrefix = "INGESTION"

type Config struct {
	BindAddress      string      `envconfig:"BIND_ADDRESS" required:"true"`
	OTelGRPCEndpoint string      `envconfig:"OTEL_GRPC_ENDPOINT" required:"true"`
	OTelServiceName  string      `envconfig:"OTEL_SERVICE_NAME" required:"true" default:"ingestion-service"`
	Kafka            KafkaConfig `envconfig:"KAFKA"`
}

type KafkaConfig struct {
	Brokers []string `envconfig:"SEED_BROKERS" required:"true"`
}

func LoadConfig(envFiles ...string) (*Config, error) {
	var cfg Config
	if err := godotenv.Load(envFiles...); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	if err := envconfig.Process(ConfigPrefix, &cfg); err != nil {
		return nil, fmt.Errorf("failed to process env variables: %w", err)
	}

	return &cfg, nil
}
