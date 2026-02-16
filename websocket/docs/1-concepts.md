# WebSocket

WebSocket is a computer communications protocol, providing a simultaneous two-way communication channel over a single Transmission Control Protocol (TCP) connection. It shifts cost from per-interaction overhead to connection lifetime and state management.

## 1. Protocol

WebSocket can operate over HTTP ports and is designed to work with HTTP proxies.

Although HTTP persistent connections and WebSocket both reuse a single TCP connection, HTTP introduces latency due to:

- Server routing and request lifecycle overhead
- Repeated header transmission with associated parsing and serialization costs per request
- Response blocking, the server cannot send independent messages until the current response completes

## 2. Connection Lifecycle

The WebSocket connection lifecycle generally consists of four primary phases or states: the Opening Handshake, Data Transfer, Closing Handshake, and final Closure.

### 2.1. Opening Handshake

State: `CONNECTING`

The client initiates the connection by sending a standard **HTTP** GET request to the server with specific headers, including `Upgrade: websocket` and `Connection: Upgrade`

- If the server supports the protocol, it responds with an HTTP 101 Switching Protocols status code and corresponding headers (Upgrade, Connection, Sec-WebSocket-Accept) to confirm the upgrade.

### 2.2. Data Transfer

State: `OPEN` or `CLOSING`

Once the handshake is successful, the connection is upgraded from HTTP to a persistent, full-duplex WebSocket connection over the same underlying TCP connection.

- Both the client and server can send messages (data frames, ping/pong control frames) back and forth independently and in real-time.
- Either the client or the server can initiate termination by sending a "close frame" message, which includes a status code indicating the reason for closure. The receiving party then responds with its own close frame to acknowledge the request.

### 2.3. Connection Closure

State: `CLOSED`

After both parties have sent and received the closing frames, the underlying TCP connection is terminated.

## 3. Limits in Practice

- **Client-side TCP Port Limit**: A single server listening on a single port (only one fixed IP) can theoretically handle up to 65,535 (16-bit) _concurrent_ connections from distinct client IP addresses, as this is the maximum number of available TCP ports on the client side.
- **Server-side File Descriptor Limit**: Every open socket consumes a file descriptor. On Linux/macOS, the default soft limit (ulimit -n) is often 1024 and the recommend hard limit is 65,535. There is effectively no fixed port-based limit on the server side for WebSocket concurrency.
- **Memory Usage**: Each WebSocket connection needs memory for TCP buffers, goroutine stack (~2 KB minimum, grows if needed), connection state, etc.
