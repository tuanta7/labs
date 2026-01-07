package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

func handleStats(w http.ResponseWriter, r *http.Request) {
	m := getMemStats()

	stats := map[string]interface{}{
		"heap_alloc_kb":   bToKb(m.HeapAlloc),
		"heap_alloc_mb":   bToMb(m.HeapAlloc),
		"total_alloc_mb":  bToMb(m.TotalAlloc),
		"sys_mb":          bToMb(m.Sys),
		"num_gc":          m.NumGC,
		"heap_objects":    m.HeapObjects,
		"goroutines":      runtime.NumGoroutine(),
		"garbage_size_mb": len(garbage),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleAllocate(w http.ResponseWriter, r *http.Request) {
	sizeMB := 10 // Default 10 MB
	allocateGarbage(sizeMB)

	m := getMemStats()
	response := map[string]interface{}{
		"message":       fmt.Sprintf("Allocated %d MB of garbage", sizeMB),
		"heap_alloc_mb": bToMb(m.HeapAlloc),
		"garbage_size":  len(garbage),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleRelease(w http.ResponseWriter, r *http.Request) {
	beforeStats := getMemStats()
	releaseGarbage()
	afterStats := getMemStats()

	response := map[string]interface{}{
		"message":        "Released garbage references",
		"before_heap_mb": bToMb(beforeStats.HeapAlloc),
		"after_heap_mb":  bToMb(afterStats.HeapAlloc),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGC(w http.ResponseWriter, r *http.Request) {
	stats := triggerGC()

	response := map[string]interface{}{
		"message":        "GC triggered",
		"timestamp":      stats.Timestamp,
		"before_heap_kb": bToKb(stats.BeforeHeapAlloc),
		"after_heap_kb":  bToKb(stats.AfterHeapAlloc),
		"cleaned_kb":     bToKb(stats.Cleaned),
		"num_gc":         stats.NumGC,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
	historyLock.RLock()
	defer historyLock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gcHistory)
}
