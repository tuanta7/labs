# Stateful

## 1. TCP Failure Detection

A TCP connection is considered persisted not through explicit coordination between machines, but through the continued validity of protocol state maintained independently at both endpoints.

### 1.1. Application-layer mitigation

For idle connections, TCP does not continuously exchange messages by default. The connection remains valid indefinitely unless one side sends a FIN or RST, or the network path fails and the failure is detected.

Because TCP failure detection can be slow, protocols frequently implement heartbeats.

- WebSocket uses ping/pong frames for this purpose.
- Gateways such as NGINX rely on read/write timeouts and socket errors.

### 1.2. HTTP Failure Detection for TCP?

While built on top of TCP – a stateful protocol, HTTP is stateless itself.

- HTTP has **NO** intrinsic mechanism for detecting abrupt peer failure or half-open connections.
- The term `keep-alive` in HTTP/1.1 refers to connection reuse, not liveness probing.

## 2. NGINX Layer 7 Proxying

Reference: [WebSocket Proxying](https://nginx.org/en/docs/http/websocket.html?)

The client establishes a TCP connection to the API gateway or load balancer. This is the only TCP connection visible to the client. Persistence is achieved through connection anchoring at the gateway layer, rather than through a direct end-to-end TCP connection between the client and the backend.

During the HTTP upgrade to WebSocket, the gateway selects a backend instance and opens (or reuses) a separate TCP connection to that backend. At this point, a connection pair is formed:

- Client ↔ Gateway (TCP)
- Gateway ↔ Backend (TCP)

After the upgrade completes, the gateway switches from HTTP request routing to stream proxy mode (**tunneling**). Bytes received on one TCP connection are forwarded directly to the other TCP connection with minimal inspection.

```scss
Kernel TCP Stack
 ├─ FD 12  → Client Socket (A:50001 → N:443)
 └─ FD 37  → Upstream Socket (N:41001 → S:8080)

NGINX (user space)
 ┌───────────────────────────┐
 │ Client Conn Object (C)    │───┐
 └───────────────────────────┘   │
                                 │  fixed pointer
 ┌───────────────────────────┐   │
 │ Upstream Conn Object (U)  │◄──┘
 └───────────────────────────┘
```

### Scaling

When NGINX needs to connect to a backend service (for HTTP, WebSocket, or any TCP-based proxying), it performs a system call equivalent to:

```c
connect(fd, backend_ip, backend_port)
```

The operating system selects a free port from the ephemeral port range (typically 1024-65535) and binds the socket implicitly before the connection is established. Each connection pair is represented by 2 tupes

```sh
# client_source/client_port and backend_source/backend_port are fixed values

1. (client_source, client_port, nginx_source, nginx_port)
2. (nginx_source, nginx_port, backend_source, backend_port)
```

For example, having an NGINX node with ~ 50k ephemeral ports and 3 backend nodes, there will be ~ 150k concurrent upstream TCP connections
