package notification

import (
	"context"
	"fmt"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/amqp"
)

type Consumer struct {
	rabbitmq amqp.Consumer
	handlers map[string]amqp.ConsumerHandler
}

func NewConsumer(rabbitmq amqp.Consumer) *Consumer {
	return &Consumer{
		rabbitmq: rabbitmq,
	}
}

func (c *Consumer) RegisterHandler(queue string, handler amqp.ConsumerHandler) {
	if c.handlers == nil {
		c.handlers = make(map[string]amqp.ConsumerHandler)
	}

	c.handlers[queue] = handler
}

func (c *Consumer) ConsumePushNotificationQueue(ctx context.Context) error {
	return c.consume(ctx, domain.PushNotificationQueue)
}

func (c *Consumer) consume(ctx context.Context, queue string) error {
	h, ok := c.handlers[queue]
	if !ok {
		return fmt.Errorf("handler not found for queue %s", domain.PushNotificationQueue)
	}

	return c.rabbitmq.Consume(ctx, queue, "", false, false, h)
}
