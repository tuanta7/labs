# Socket.io

## 1. Supported Protocols

- WebSocket
- Long Polling

## 2. Sticky Session

Sticky session, also known as session affinity, is a load balancing strategy where a client consistently connects to the same server across a session or multiple reconnects.

- Native WebSocket connections are stateful and remain bound to a single server instance for the lifetime of the connection; therefore, load balancer–level sticky sessions are not inherently required for correct WebSocket operation, provided that the connection itself is not redistributed mid-session.
- Socket.IO may require sticky sessions when operating in transport fallback modes (for example, HTTP long polling). In such cases, multiple sequential HTTP requests belonging to the same logical session (typically include in-memory state logic) must be routed to the same server instance to preserve session state and ensure correct behavior.

> [!NOTE]
> Some clients, due to restrictive firewalls, proxies, or legacy environments, won't be able to establish a WebSocket connection at all.
