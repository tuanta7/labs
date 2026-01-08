# Go Garbage Collector

Reference: [go.dev](https://go.dev/doc/gc-guide)

At a high level, a garbage collector (or GC, for short) is a system that recycles memory on behalf of the application by identifying which parts of memory are no longer needed. 

- **Object**: An object is a dynamically allocated piece of memory that contains one or more Go values.
- **Pointer**: A memory address that references any value within an object.
- **Object Graph**: A graph of objects and pointers to other objects.
- **Scanning**: The process of walking the object graph.
- **Root**: The starting points from which the garbage collector determines which objects in memory are still **reachable** and therefore should not be collected

## 1. The GC cycle

Go GC is a mark-sweep GC, it broadly operates in two phases: the mark phase, and the sweep phase (starting with sweeping, turning off, then marking). It's not possible to release memory back to be allocated until all memory has been traced, because there may still be an un-scanned pointer keeping an object alive.

### 1.1. Tracing Garbage Collection

Reference: [leapcell.io](https://leapcell.io/blog/understanding-memory-management-in-go)

GC roots are defined as all memory references that are guaranteed to remain reachable regardless of heap state.

- To identify live memory, the GC walks the object graph starting at the program's roots - pointers that identify objects that are definitely in-use by the program. Two examples of roots are local variables and global variables.

### 1.2. GOGC (variable, can be % or "off")

GOGC is an environment variable, which defaults to 100(%). At a high level, GOGC determines the trade-off between GC CPU and memory.

- GC is triggered when the live heap size is larger than the target heap memory (double when GOGC=100)
- Lower GOGC values cause the garbage collector to run more frequently, as the trigger point for new allocations is lower. As GOGC decreases, the peak memory requirement decreases at the expense of additional CPU overhead.
- At GOGC=0, the GC will still run when memory hit the limit set by the `GOMEMLIMIT` environment variable

The target heap memory is defined as:

$$
Target Heap Memory = Live Heap + (Live Heap + GC Roots) * GOGC / 100
$$

Live heap memory is memory that was determined to be live by the previous GC cycle, while new heap memory is any memory allocated in the current cycle, which may or may not be live by the end. 

> [!NOTE]
> As an example, consider a Go program with a live heap size of 8 MiB, 1 MiB of goroutine stacks, and 1 MiB of pointers in global variables. 
> - With a GOGC value of 100, the amount of new memory that will be allocated before the next GC runs will be 10 MiB, or 100% of the 10 MiB of work, for a total heap footprint of 18 MiB (8 + 10). 
> - With a GOGC value of 50, then it'll be 50%, or 5 MiB, total heap footprint = 13 MiB. 
> - With a GOGC value of 200, it'll be 200%, or 20 MiB, total heap footprint = 28 MiB.

### 1.3. GC Trace Lines

Set `GODEBUG=gctrace=1` to track every time GC runs. A typical trace line provides detailed metrics about each collection cycle

```shell
GODEBUG=gctrace=1 go run main.go

gc 1 @0.191s 0%: 0.16+6.5+0.079 ms clock, 1.3+0.32/1.1/0+0.63 ms cpu, 3->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.222s 0%: 0.025+0.73+0.64 ms clock, 0.20+0.053/0.64/0+5.1 ms cpu, 3->3->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 3 @0.271s 0%: 0.086+1.1+0.004 ms clock, 0.68+0.38/1.2/0+0.036 ms cpu, 3->3->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
...
GC forced 
gc 15 @120.664s 0%: 0.18+1.1+0.005 ms clock, 1.4+0/2.1/0+0.045 ms cpu, 2->2->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
```

- `gc 1`, `gc 2` indicates the sequential GC cycle number since program start.
- `@0.191s`, `@120.664s` indicates the wall-clock time elapsed since program start when the GC cycle completed.
- `n%` represents the fraction of total CPU time spent in GC during the interval since the previous GC cycle.

See [forcegcperiod](https://github.com/golang/go/blob/master/src/runtime/proc.go#L6265) for the maximum time in nanoseconds between garbage collecting cycles.

```go
// go1.25
var forcegcperiod int64 = 2 * 60 * 1e9 // 2 minutes
```

## 2. Virtual Memory

Reference: [200lab.io](https://200lab.io/blog/golang-cap-phat-bo-nho-nhu-the-nao)

Virtual memory is an abstraction over physical memory provided by the operating system to isolate programs from one another. It's also typically acceptable for programs to reserve virtual address space that doesn't map to any physical addresses at all.