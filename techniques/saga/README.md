# EV Charging

A saga orchestrator using RabbitMQ as the messaging queue with full acknowledgment and QoS, designed for handling reservations to charge electric vehicles.

## Use Case: Start Charging Session

A driver attempts to start charging at a station through a mobile application. The services involved are:

- **Reservation Service**: Reserve a charging connector at a specific station.
- **Payment Service**: Authorize payment for the charging session.
- **Charging Device Service**: Activate the charger hardware.

### Saga States

| Status             | Description                                |
|--------------------|--------------------------------------------|
| `STARTED`          | Saga created, first command about to send  |
| `RESERVE_PENDING`  | Waiting for reservation service response   |
| `RESERVE_OK`       | Reservation confirmed                      |
| `PAYMENT_PENDING`  | Waiting for payment service response       |
| `PAYMENT_OK`       | Payment authorized                         |
| `CHARGING_PENDING` | Waiting for charging service response      |
| `COMPLETED`        | All steps succeeded                        |
| `COMPENSATING`     | Rolling back previously completed steps    |
| `FAILED`           | Saga failed after compensation             |

## Forward Transaction Flow

1. **Reserve** → Orchestrator publishes `EXECUTE` command to `reservation` queue
2. Reserve service processes, replies with success → Orchestrator advances
3. **Payment** → Orchestrator publishes `EXECUTE` command to `payment` queue
4. Payment service authorizes, replies with success → Orchestrator advances
5. **Charging** → Orchestrator publishes `EXECUTE` command to `charging` queue
6. Charging service activates hardware, replies with success → Saga **COMPLETED**

## Failure & Compensation Flow

If any step fails, the orchestrator publishes `ROLLBACK` commands to all previously completed steps **in reverse order**:

- Payment fails → Orchestrator sends `ROLLBACK` to reservation service
- Charging fails → Orchestrator sends `ROLLBACK` to payment, then reservation

### RabbitMQ Guarantees

- **QoS**: `prefetchCount=1` on every consumer — one message at a time
- **Manual Ack**: `autoAck=false` everywhere; messages are `Ack`'d only after successful processing
- **Nack + Requeue**: On transient errors, messages are `Nack`'d with `requeue=true`
- **Persistent Messages**: `DeliveryMode: 2` (persistent) on all published messages
- **Durable Queues**: All queues declared with `durable: true`


