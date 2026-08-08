# Goroutine Leak

A goroutine leak occurs when goroutines are created but never terminate, accumulating in memory over time.

Unlike memory leaks from objects, goroutine leaks are about "processes" that never exit, they keep running indefinitely, consuming resources and preventing garbage collection of associated data.

## Detection Methods

### Metrics

Add `runtime.NumGoroutine()` to Prometheus as a custom metric. In a healthy system, the number of goroutines should return to baseline after traffic decreases.

```txt
# TYPE go_goroutines gauge
go_goroutines 37
```

Compare against request rate (HTTP RPS, Kafka consumer throughput, Queue depth, etc.), if traffic drops but goroutines keep increasing, that strongly indicates leaked workers.

### Profiling

Identifying which goroutines are leaking generally requires stack traces from the goroutine profile (pprof) or equivalent runtime diagnostics. Metrics alone can strongly suggest a leak but rarely pinpoint the source.

```sh
# fetches the profile over HTTP
go tool pprof http://host/debug/pprof/goroutine
```

## Common Patterns

### Unbuffered Channel with no receiver

```go
ch := make(chan int)  // No buffer, capacity 0

go func() {
  result := expensiveComputation()
  ch <- result  // BLOCKS
}()
```

Ensure a receiver exists or use a buffered channel

```go
ch := make(chan int, 10)

go func() {
  result := expensiveComputation()
  ch <- result  // Writes to buffer and returns immediately
}()
```

### Missing context cancellation
