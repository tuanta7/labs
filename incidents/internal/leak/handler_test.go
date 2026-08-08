package main

import (
	"runtime"
	"testing"
)

func BenchmarkOperationWithLeakDetection(b *testing.B) {
	before := runtime.NumGoroutine()

	for i := 0; i < b.N; i++ {
		Handler() // Runs hundreds of times
	}

	runtime.GC()
	after := runtime.NumGoroutine()
	if after > before {
		b.Fatalf("Leak after %d iterations: %d -> %d", b.N, before, after)
	}
}
