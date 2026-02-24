# Kafka Failover

## 1. Consumer Rebalancing

Reference: [Kafka Rebalancing Explained](https://www.confluent.io/learn/kafka-rebalancing/)

Rebalancing is triggered by membership changes, session timeouts, and exceeding the maximum poll interval.

## 2. KRaft (Kafka-Raft)

In KRaft mode each Kafka server can be configured as a controller, a broker, or both using the `process.roles` property.
