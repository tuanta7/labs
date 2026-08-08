# Memory Leak

Memory growth alone does not necessarily indicate a memory leak. Go's runtime deliberately grows the heap and may not immediately return memory to the OS. The goal is to distinguish normal heap growth from memory that remains unexpectedly live.

When a Go service runs out of memory (OOM), the runtime will attempt to trigger garbage collection to free memory

```shell

```

## Detection Methods

### Metrics

### Profiling: Capture heap profiles

## Common Patterns

