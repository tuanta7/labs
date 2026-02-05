package amqp

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
}

func NewClient(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &Client{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *Client) Close() error {
	if c.channel != nil {
		_ = c.channel.Close()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func (c *Client) QueueDeclare(queue string, durable, autoDelete, exclusive bool, args Arguments) (Queue, error) {
	q, err := c.channel.QueueDeclare(queue, durable, autoDelete, exclusive, false, amqp.Table(args))
	if err != nil {
		return Queue{}, err
	}

	return Queue{
		Name:      q.Name,
		Consumers: q.Consumers,
		Messages:  q.Messages,
	}, nil
}

func (c *Client) ExchangeDeclare(exchange, kind string, durable, autoDelete bool, args Arguments) error {
	return c.channel.ExchangeDeclare(exchange, kind, durable, autoDelete, false, false, amqp.Table(args))
}

func (c *Client) QueueBind(queue, exchange, key string, args Arguments) error {
	return c.channel.QueueBind(queue, exchange, key, false, amqp.Table(args))
}

func (c *Client) ExchangeBind(destination, source, key string, args Arguments) error {
	return c.channel.ExchangeBind(destination, source, key, false, amqp.Table(args))
}

func (c *Client) Publish(ctx context.Context, exchange, key string, mandatory bool, body []byte) error {
	return c.channel.PublishWithContext(ctx, exchange, key, mandatory, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (c *Client) Consume(ctx context.Context, queue, consumer string, autoAck, exclusive bool, handler ConsumerHandler) error {
	deliveryCh, err := c.channel.Consume(queue, consumer, autoAck, exclusive, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range deliveryCh {
		err = handler.Handle(ctx, msg.Body)
		if err != nil {
			_ = msg.Nack(false, false)
			return err
		}

		_ = msg.Ack(false)
	}

	return nil
}
