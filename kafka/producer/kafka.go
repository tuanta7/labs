package main

import "github.com/twmb/franz-go/pkg/kgo"

type KafkaProducer struct {
	client *kgo.Client
}
