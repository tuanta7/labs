# Routing

Routing in AMQP 0-9-1 is based on the concept of **exchanges**, **queues**, and **bindings**. Producers publish messages to exchanges, and exchanges use bindings (with optional routing keys) to decide how to deliver messages to queues or other exchanges. Consumers then consume from queues.

## 1. Declaring

Queues and exchanges must exist before they can be used. Declaring them in application code is a common approach, but not mandatory if a system-level setup step already created them.

- Declaring a queue/exchange that already exists with the same attributes is safe.
- Declaring with conflicting attributes causes a channel-level error (code 406, _PRECONDITION_FAILED_).

### 1.1. Queue

```go
q, err := ch.QueueDeclare(
  "hello", // name: identifier for the queue
  false,   // durable: survive broker restart
  false,   // auto-delete: delete when unused
  false,   // exclusive: only accessible to this connection
  false,   // no-wait: don't wait for server response
  nil,     // arguments: extra features (e.g., TTL, dead-lettering)
)
if err != nil {
    panic(err)
}
```

### 1.2. Exchange

Declaring an exchange explicitly defines the routing type and durability

- Declaring an exchange with an empty string refers to the default exchange, a built-in direct exchange that routes based on queue name.
- Well-known system exchanges like `amq.direct`, `amq.topic`, or `amq.fanout` are usually pre-declared and can be used without `ExchangeDeclare`.

```go
err := ch.ExchangeDeclare(
    "exchange-name", // name: identifier for the exchange
    "fanout",        // type: direct | fanout | topic | headers
    true,            // durable: survive broker restart
    false,           // auto-delete: delete when no queues are bound
    false,           // internal: exchange not directly publishable by clients
    false,           // no-wait: don't wait for server response
    nil,             // arguments: extra features (e.g., alternate-exchange)
)
if err != nil {
    panic(err)
}
```

#### Exchange Types

| Exchange type | Queue to be routed to                                        |
| ------------- | ------------------------------------------------------------ |
| direct        | exact `routingKey` match                                     |
| fanout        | broadcast to all bound queues                                |
| topic         | wildcard matching with `*` (single word) and `#`(multi-word) |
| headers       | match based on message headers                               |

## 2. Binding

A binding links a queue (or exchange) to an exchange and optionally filters messages via a routing key.

- Without at least one binding, messages sent to a non-default exchange will be dropped.
- Bindings are usually created on the consumer side (or in a bootstrap step) to ensure the correct routing topology exists before consumers begin. Producers should not typically perform `QueueBind` — this keeps publishing logic lightweight and avoids accidental topology coupling.

Queue binding: This ensures that any message published to `exchangeName` with a matching `routingKey` will be routed to `queueName`.

```go
err := ch.QueueBind(
    queueName,    // name of the queue to bind
    routingKey,   // key to match messages against
    exchangeName, // name of the exchange to bind to
    false,        // no-wait: don't wait for server response
    nil,          // arguments: extra binding features
)
if err != nil {
    panic(err)
}
```

Exchange-to-exchange binding: Messages published to `sourceExchange` are routed to `destinationExchange` as if they were directly published there.

```go
err := ch.ExchangeBind(
    destinationExchange, // exchange to route messages into
    routingKey,          // key to match messages against
    sourceExchange,      // exchange to route messages from
    false,               // no-wait: don't wait for server response
    nil,                 // arguments: extra binding features
)
if err != nil {
    panic(err)
}
```

## 3. Putting It Together

A typical consumer setup may look like this

```go
// Declare exchange
err := ch.ExchangeDeclare("exchange-name", "direct", true, false, false, false, nil);
if err != nil {
    panic(err)
}

// Declare queue
q, err := ch.QueueDeclare("queue-name", true, false, false, false, nil)
if err != nil {
    panic(err)
}

// Bind queue to exchange with routing key
if err := ch.QueueBind(q.Name, "routing-key", "exchange-name", false, nil); err != nil {
    panic(err)
}
```

And a publisher might only need

```go
err := ch.Publish(
    "exchange-name", // exchange: target exchange for the message
    "routing-key",   // routing key: determines message destination
    false,           // mandatory: return undeliverable messages to sender
    false,           // immediate: deprecated, unused
    amqp.Publishing{
        ContentType: "text/plain",
        Body:        []byte("hello world"),
    },
)
if err != nil {
    panic(err)
}
```
