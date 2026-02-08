## AMQP Go Client

The official AMQP 0-9-1 Go client is `rabbitmq/amqp091-go` (forked from `streadway/amqp`), which is maintained by the RabbitMQ team. In this library, a message sent from publisher is called a `Publishing` and a message received to a consumer is called a `Delivery`.

```sh
go get github.com/rabbitmq/amqp091-go
```

## 1. Connection & Channels

A single connection is intended to last for the full lifetime of the process.

```go
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
if err != nil {
    panic(err)
}
defer conn.Close()
```

Channels are created per logical task (publisher, consumer group, etc.) and are not thread‑safe.

```go
ch, err := conn.Channel()
if err != nil {
    panic(err)
}
defer ch.Close()
```

## 2. Queues

A queue is a sequential data structure with two primary operations: an item can be enqueued (added) at the tail and dequeued (consumed) from the head

### 2.1. Properties

Queues have properties that define how they behave. There is a set of mandatory properties and a map of optional ones

- **Name**: The name of the queue. In AMQP 0-9-1, the broker can generate a unique queue name on behalf of an app
- **Durability**: Queues can be durable or transient. Metadata of a durable queue is stored on disk, while metadata of a transient queue is stored in memory when possible. 
- **Exclusive**: Queue that used by only one connection and will be deleted when that connection closes
- **Auto-delete**: Queue that has had at least one consumer is deleted when last consumer unsubscribes
- **Arguments** (optional): Used by plugins and broker-specific features such as message TTL, queue length limit, etc.

> [!NOTE]
> When an auto-delete queue is created, and there are no consumers declared on it, the queue will not be automatically deleted because the condition of "last consumer unsubscribes" for deletion is never met.

### 2.2. Quorum Queues

## 3. Exchanges

In AMQP 0-9-1, exchanges are the entities where publishers publish messages that are then routed to a set of queues or streams.

### 3.1. Properties

Exchanges have several key properties that can be specified at declaration times

- **Durability**: Just like queues, exchanges can be durable or transient. However, transient exchanges are very rarely used in practice.
- **Auto-delete**: Auto-deleted exchanges are deleted when their last binding is removed.
- **Arguments** (optional): Optional exchange arguments, also known as "x-arguments" because of their field name in the AMQP 0-9-1 protocol, is a map (dictionary) of arbitrary key/value pairs that can be provided by clients when an exchange is declared.

### 3.2. Dead Letter Exchange

Messages from a queue can be "dead-lettered", which means these messages are republished to an exchange when any of the following 4 events occur