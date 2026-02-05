# RabbitMQ

RabbitMQ is a message broker, it accepts and forwards messages.

## 1. Message Lifecycle

When a message enters RabbitMQ it passes through three distinct phases

| Phase     | Internal state | Description                                                                                                 |
| --------- | -------------- | ----------------------------------------------------------------------------------------------------------- |
| Queued    | Ready          | The message resides on disk or in memory, assigned to a queue but not yet delivered.                        |
| In‑flight | Unacknowledged | The message has been sent to a consumer, yet the broker has not received an explicit acknowledgment for it. |
| Completed | Acked / Nacked | The consumer has responded, the broker then removes or re‑queues the message accordingly.                   |

Consumer acknowledgments (delivery‑side reliability)

| Method                     | Effect                                                                    | Typical use        |
| -------------------------- | ------------------------------------------------------------------------- | ------------------ |
| `basic.ack`                | Confirms successful processing; message is deleted.                       | Normal completion. |
| `basic.nack requeue=true`  | Signals failure; message is re‑queued for another attempt.                | Transient errors.  |
| `basic.nack requeue=false` | Routes the message to a dead‑letter exchange (if configured) or drops it. | Poison messages.   |
| `basic.reject`             | Single‑message variant of NACK; identical semantics.                      | Edge cases.        |

With autoAck=true the broker discards the message immediately after sending it, eliminating the unacknowledged phase but forfeiting at‑least‑once guarantees.

> [!CAUTION]
> Publisher confirms are an independent protocol feature

| Direction         | Mechanism               | Guarantees                                                              |
| ----------------- | ----------------------- | ----------------------------------------------------------------------- |
| Producer → Broker | Publisher Confirm       | Message has been persisted and routed to at least one queue.            |
| Broker → Consumer | Delivery Acknowledgment | Message has been processed (or intentionally rejected) by the consumer. |

## 2. Channels

Most applications need multiple logical connections to the broker. However, it is undesirable to keep many TCP connections open at the same time because doing so consumes system resources and makes it more difficult to configure firewalls.

- AMQP 0-9-1 connections are multiplexed with channels.
- Channels can be thought of as lightweight connections that share a single TCP connection.
- Much like connections, channels are meant to be long lived (using channel pool).
- The client cannot be configured to allow for more channels than the server configured maximum

## 3. Queues & Exchanges

A queue in RabbitMQ is an ordered collection of messages. Messages are enqueued and dequeued (delivered to consumers) in a FIFO manner.

> [!NOTE]
> Clients publish to Exchanges, not to Queues.

In AMQP 0-9-1, exchanges are the entities where publishers publish messages that are then routed to a set of queues or streams.

- Exchanges routes all messages that flow through them to one or more queues, streams, or other exchanges.
- Every exchange belongs to one virtual host (logical groups of entities)

The routing between producers and consumer queues is via Bindings. These bindings form the logical topology of the broker.

### Streams

RabbitMQ Streams is a persistent replicated data structure that can complete the same tasks as queues: they buffer messages from producers that are read by consumers.

## References

- [RabbitMQ | Go Tutorials](https://www.rabbitmq.com/tutorials/tutorial-one-go)
- [GitHub | RabbitMQ Go Tutorials](https://github.com/rabbitmq/rabbitmq-tutorials/tree/main/go)
- [Acknowledger Interface](https://github.com/rabbitmq/amqp091-go/blob/v1.10.0/delivery.go#L19)
