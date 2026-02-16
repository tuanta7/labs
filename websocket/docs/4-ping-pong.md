# WebSocket Ping/Pong Implementation

## Server vs Client: Which Should Send Ping?

### ✅ **Most Common: Server Sends PING**

The **server** typically sends PING frames, and the **client** automatically responds with PONG frames.

## Why Server-Side Ping is Standard

### 1. **Connection Management**
- Server needs to detect dead/stale connections
- Free up server resources (memory, file descriptors)
- Clean up connection pools

### 2. **Scalability**
- One server manages thousands of clients
- Server needs to know which connections are alive
- Prevents resource exhaustion

### 3. **Protocol Design**
- WebSocket protocol was designed for server-initiated pings
- Browsers automatically respond to server pings with pongs
- No JavaScript needed on client side

### 4. **Network Reality**
- Clients can have flaky connections (mobile, wifi)
- Clients might be behind NAT/firewalls
- Server can enforce connection health policy

## Current Implementation

### Server Side (handler.go)

```go
const (
    pongWait   = 60 * time.Second        // Wait for pong from client
    pingPeriod = (pongWait * 9) / 10     // Send ping every 54 seconds
)
```

**Flow:**
1. Server sends PING frame every 54 seconds
2. Client browser automatically responds with PONG
3. Server has 60 second deadline to receive PONG
4. If no PONG received, connection is considered dead

### Client Side (HTML)

The HTML client has optional "application-level" ping/pong:
- Sends text message "PING" every 5 seconds
- This is NOT the WebSocket protocol ping/pong
- This is for testing/demonstration purposes

## Two Types of Ping/Pong

### 1. **Protocol-Level** (Recommended)
```
Server → PING frame (opcode 0x9)
Client → PONG frame (opcode 0xA)
```
- Handled by WebSocket library
- Transparent to application
- Browser automatically responds
- Used for connection health

### 2. **Application-Level** (Optional)
```
Client → "PING" text message
Server → "PONG" text message
```
- Custom business logic
- Requires JavaScript code
- For testing or specific needs
- Visible in message logs

## Best Practices

### ✅ DO:
- Implement server-side protocol ping/pong
- Set reasonable timeouts (60s is common)
- Log ping/pong for debugging
- Close connections that don't respond

### ❌ DON'T:
- Don't rely only on client-side pings
- Don't set ping intervals too short (< 30s)
- Don't ignore pong timeout errors
- Don't mix up protocol vs application pings

## Testing

### With Browser DevTools:
1. Open Network tab → WS filter
2. Look for "Ping" and "Pong" frames
3. Check frame types (not text messages)

### With Your HTML Client:
1. Enable "Enable Ping/Pong" checkbox
2. This sends application-level pings
3. See "PING" in message log
4. Protocol pings are invisible in logs

## Current Setup Summary

| Feature | Implementation | Type |
|---------|---------------|------|
| Server → PING | ✅ Every 54s | Protocol |
| Client → PONG | ✅ Automatic | Protocol |
| Client → "PING" | ✅ Optional | Application |
| Server → Echo | ✅ Yes | Application |

Your implementation now follows WebSocket best practices with server-side ping/pong! 🎉

