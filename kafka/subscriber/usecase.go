package main

import "context"

type Consumer interface {
	Consume(ctx context.Context, topic string, groupID string, handler func(key, value []byte) error) error
}

type ConsumerUseCase struct {
	consumer Consumer
}

func NewConsumerUseCase(consumer Consumer) *ConsumerUseCase {
	return &ConsumerUseCase{
		consumer: consumer,
	}
}

func (c *ConsumerUseCase) ConsumeLocation(ctx context.Context) {}
