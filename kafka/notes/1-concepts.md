# Apache Kafka

Reference: [Apache Kafka 4.1](https://kafka.apache.org/41/)

Kafka is a distributed system consisting of servers and clients that communicate via a high-performance TCP network protocol.

## 1. Concepts

Reference: [Kafka Introduction](https://kafka.apache.org/41/getting-started/introduction/)

An event records the fact that something happened in the world or in your business. It is also called record or message. When you read or write data to Kafka, you do this in the form of events.

Conceptually, an event has a key, value, timestamp, and optional metadata headers. For example:

- **Event key**: "Alice"
- **Event value**: "Made a payment of $200 to Bob"
- **Event timestamp**: "Jun. 25, 2020 at 2:06 p.m."

### 1.1. Producers & Consumers

Producers are those client applications that publish (write) events to Kafka and consumers are those that subscribe to (read and process) these events

### 1.2. Topics

Events are organized and durably stored in topics. Very simplified, a topic is similar to a folder in a filesystem, and the events are the files in that folder.

- Events in a topic can be read as often as needed—unlike traditional messaging systems, events are not deleted after consumption.
- Kafka’s performance is effectively constant with respect to data size, so storing data for a long time is perfectly fine.

Topics are partitioned , meaning a topic is spread over a number of buckets located on different Kafka brokers.

- Events with the same event key (e.g., a customer or vehicle ID) are written to the same partition
- A single consumer may safely consume from multiple partitions.

![](../../.imgs/topic.png)

### 1.3. Consumer Group

A consumer group is a set of consumers that collectively consume messages from one or more topics. The group is identified by a `group.id`.

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
