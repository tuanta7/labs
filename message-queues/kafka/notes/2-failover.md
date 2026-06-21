# Kafka Failover

## 1. Consumer Rebalancing

Reference: [Kafka Rebalancing Explained](https://www.confluent.io/learn/kafka-rebalancing/)

Rebalancing is triggered by membership changes, session timeouts, and exceeding the maximum poll interval.

## 2. Message Delivery Guarantees

Reference: [Kafka Message Delivery Guarantees](https://docs.confluent.io/kafka/design/delivery-semantics.html#exactly-once-support)

By default Kafka guarantees at-least-once delivery. 

### Message Duplication

Kafka restarting can lead to duplicate messages primarily due to a mismatch between message processing and offset commits

## 3. KRaft

In KRaft mode each Kafka server can be configured as a controller, a broker, or both using the `process.roles` property.