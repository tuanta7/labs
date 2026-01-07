package main

import (
	"runtime"
	"sync"
	"time"
)

// GCStats holds memory statistics before and after GC
type GCStats struct {
	Timestamp       string `json:"timestamp"`
	BeforeHeapAlloc uint64 `json:"before_heap_alloc"`
	AfterHeapAlloc  uint64 `json:"after_heap_alloc"`
	Cleaned         uint64 `json:"cleaned"`
	NumGC           uint32 `json:"num_gc"`
	TotalAlloc      uint64 `json:"total_alloc"`
	Sys             uint64 `json:"sys"`
}

var (
	gcHistory   []GCStats
	historyLock sync.RWMutex
)

// triggerGC triggers garbage collection and records stats
func triggerGC() GCStats {
	beforeStats := getMemStats()
	runtime.GC()
	afterStats := getMemStats()

	cleaned := uint64(0)
	if beforeStats.HeapAlloc > afterStats.HeapAlloc {
		cleaned = beforeStats.HeapAlloc - afterStats.HeapAlloc
	}

	stats := GCStats{
		Timestamp:       time.Now().Format("15:04:05.000"),
		BeforeHeapAlloc: beforeStats.HeapAlloc,
		AfterHeapAlloc:  afterStats.HeapAlloc,
		Cleaned:         cleaned,
		NumGC:           afterStats.NumGC,
		TotalAlloc:      afterStats.TotalAlloc,
		Sys:             afterStats.Sys,
	}

	historyLock.Lock()
	gcHistory = append(gcHistory, stats)
	// Keep only the last 50 entries
	if len(gcHistory) > 50 {
		gcHistory = gcHistory[1:]
	}
	historyLock.Unlock()

	return stats
}
