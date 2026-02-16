# Uber Clone

This repository will grow gradually to keep the codebase simple; early package splitting will be avoided whenever possible.

- Accept interfaces, return structs.
- Database and queue adapters (wrappers/clients) should be implemented in the `pkg` directory.
- Define the adapter interfaces where they are consumed. These adapter interfaces should abstract driver or library differences for the same backend (for example, switching from `lib/pq` to `pgx`, or `go-redis` to `redigo`), not high‑level changes such as PostgreSQL to MongoDB or RabbitMQ to Kafka.
- Sometimes a library can be trusted to be never deprecated, so it can be used directly without mapping all of its types to a local one.
- The interfaces that abstract the backend switching (such as PostgreSQL to MongoDB) should be defined in the higher level, like the `repository` layer.

## 1. High Level Design

Reference: [System Design School | Design Uber, Lyft](https://systemdesignschool.io/problems/uber/solution)

## 2. WebSocket

Reference: [Scaling WebSockets](https://ably.com/topic/the-challenge-of-scaling-websockets)

WebSocket services are used to maintain a live feed of driver locations for both drivers (update their location) and passengers (get the driver updates). In real-world production like Uber, QUIC/HTTP3 - a more modern technology is used instead.

> [!NOTE]
> Some clients, due to restrictive firewalls, proxies, or legacy environments, won't be able to establish a WebSocket connection at all.

## 3. ClickHouse

Trip History Storage

A time-series/column-oriented database like TimescaleDB (PostgreSQL extension) or ClickHouse is recommended. This storage pattern allows efficient historical queries for analytics or trip reconstruction.

## 4. Redis

Latest Location Storage with Geospatial Indexing

A high-performance key-value store is preferred because only the latest update is needed, overwriting the record is sufficient.

## 5. Kafka

Messaging Layer

A location-tracking workload behaves like a high-frequency telemetry stream. Kafka's partitioned log design allows location updates from thousands of drivers can be processed without pressure on the broker.
