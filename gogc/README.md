# Go Garbage Collector

## GC Trace Lines

A typical trace line provides detailed metrics about each collection cycle

```shell
gc 1 @14.761s 0%: 0.072+0.48+0.053 ms clock, 0.87+0/0.40/0+0.63 ms cpu, 4->4->4 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 12 P
gc 2 @14.769s 0%: 0.017+0.23+0.022 ms clock, 0.21+0/0.34/0.031+0.26 ms cpu, 9->9->9 MB, 9 MB goal, 0 MB stacks, 0 MB globals, 12 P
```