package amqp

import "context"

type Arguments map[string]any

// Queue represents a snapshot of a queue
type Queue struct {
	Name      string
	Consumers int
	Messages  int
}

type Broker interface {
	QueueDeclare(queue string, durable, autoDelete, exclusive bool, args Arguments) (Queue, error)
	ExchangeDeclare(exchange, kind string, durable, autoDelete bool, args Arguments) error
	QueueBind(queue, exchange, key string, args Arguments) error
	ExchangeBind(destination, source, key string, args Arguments) error
	Close() error
}

type Publisher interface {
	Publish(ctx context.Context, exchange, key string, mandatory bool, body []byte) error
}

type ConsumerHandler interface {
	Handle(ctx context.Context, msg []byte) error
}

type Consumer interface {
	Consume(ctx context.Context, queue, consumer string, autoAck, exclusive bool, handler ConsumerHandler) error
}
