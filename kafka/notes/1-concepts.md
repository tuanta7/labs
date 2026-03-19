# Apache Kafka

Reference: [Apache Kafka 4.1](https://kafka.apache.org/41/)

Kafka is a distributed system consisting of servers and clients that communicate via a high-performance TCP network protocol.

It solves the problem of high-throughput, low-latency, and scalable event streaming, while ensuring durability and the ability to replay data.

- Traditional databases are optimized for transactions and queries
- Traditional message queues are designed for task delivery

## 1. Concepts

Reference: [Kafka Introduction](https://kafka.apache.org/41/getting-started/introduction/)

An event records (aka record or message) the fact that something happened in the world or in a business. Reading and writing data to Kafka occurs in the form of events. Conceptually, an event has a key, value, timestamp, and optional metadata headers. For example:

- **Event Key**: "Alice"
- **Event Value**: "Made a payment of $200 to Bob"
- **Event Timestamp**: "Jun. 25, 2020 at 2:06 p.m."

### 1.1. Distributed Log

### 1.2. Producers & Consumers

Producers are those client applications that publish (write) events to Kafka and consumers are those that subscribe to (read and process) these events

### 1.3. Topics & Partitions

Events are organized and durably stored in topics. A topic is similar to a folder in a filesystem, and the events are the files in that folder.

- Events in a topic can be read as often as needed; unlike traditional messaging systems, events are not deleted after consumption.
- Kafka's performance is effectively constant with respect to data size, so storing data for a long time is perfectly fine.

Topics are partitioned, meaning a topic is spread over a number of buckets located on different Kafka brokers.

- Events with the same event key (e.g., a customer or vehicle ID) are written to the same partition
- A single consumer may safely consume from multiple partitions.
- Kafka does not support decreasing the number of partitions, only increasing.

> [!CAUTION]
> Kafka guarantees message ordering only within a single partition. Messages that need to remain in order must be routed to the same partition using the same key, ensuring that related events are managed correctly.
>  
> Increasing  the number of partitions change the mapping between keys and partitions, which may break ordering guarantees.

![](../../.imgs/topic.png)

### 1.4. Consumer Group

A consumer group is a set of consumers that collectively consume messages from one or more topics. The group is identified by an ID.

- Within a group, each partition of a topic is assigned to exactly one consumer, ensuring no duplication within the group.
- Multiple consumer groups can subscribe to the same topic independently, effectively enabling publish-subscribe semantics.
- When consumers join, leave, or crash, partitions are reassigned dynamically. This provides elasticity but can cause temporary pauses (rebalancing overhead).

## 2. Common Configs

### 2.1. Brokers

#### [bootstrap.servers](https://kafka.apache.org/41/configuration/kafka-connect-configs/#connectconfigs_bootstrap.servers)

A list of host/port pairs used to establish the initial connection to the Kafka cluster.

- Clients use this list to bootstrap and discover the full set of Kafka brokers.
- Not a special broker role in the cluster.

#### [node.id](https://kafka.apache.org/41/configuration/broker-configs/#brokerconfigs_node.id)

The node ID associated with the roles this process is playing when process.roles is non-empty. This is required configuration when running in KRaft mode.

#### [process.roles](https://kafka.apache.org/41/configuration/broker-configs/#brokerconfigs_process.roles)

The roles that this process plays: 'broker', 'controller', or 'broker,controller' if it is both.

#### [controller.quorum.voters](https://kafka.apache.org/41/configuration/broker-configs/#brokerconfigs_controller.quorum.voters)

Map of id/endpoint information for the set of voters in a comma-separated list of `{id}@{host}:{port}` entries. For example: `1@localhost:9092,2@localhost:9093,3@localhost:9094`

#### [listeners](https://kafka.apache.org/41/configuration/broker-configs/#brokerconfigs_listeners)

Read more: [Kafka Listeners – Explained](https://www.confluent.io/blog/kafka-listeners-explained/)

A comma-separated list of listeners and the host/IP and Kafka port to which Kafka binds to and listens for incoming connections

- If the listener name is not a security protocol, `listener.security.protocol.map` must also be set.

### 2.2. Topic/Partitions

#### [auto.create.topics.enable](https://kafka.apache.org/41/configuration/broker-configs/#brokerconfigs_auto.create.topics.enable)

Enable auto creation of topic on the server.

- Automatically create topics when a client attempts to produce to or consume from a non-existent topic.
- Commonly disabled In modern production environments to prevent accidental or misconfigured topic creation.

#### [num.partitions](https://kafka.apache.org/41/configuration/broker-configs/#brokerconfigs_num.partitions)

The default number of log partitions per topic

- Hard upper bound on parallel consumption per consumer group

### 2.3.

## 3. Kafka Streams
