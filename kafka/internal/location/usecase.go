package location

import "context"

type KafkaProducer interface {
	Produce(ctx context.Context, topic string, key, value []byte)
}

type UseCase struct {
	producer KafkaProducer
}
