package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	KafkaBrokers       []string `env:"KAFKA_BROKERS,default=localhost:9092"`
	KafkaTopicLocation string   `env:"KAFKA_TOPIC_LOCATION,default=location"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
