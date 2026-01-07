package main

import "sync"

var (
	garbage     [][]byte
	garbageLock sync.Mutex
)

// allocateGarbage creates garbage data that GC can clean
func allocateGarbage(sizeMB int) {
	garbageLock.Lock()
	defer garbageLock.Unlock()

	// Allocate sizeMB worth of data
	for i := 0; i < sizeMB; i++ {
		buf := make([]byte, 1024*1024) // 1 MB chunks
		for j := range buf {
			buf[j] = byte(j % 256)
		}
		garbage = append(garbage, buf)
	}
}

// releaseGarbage releases the allocated garbage so GC can clean it
func releaseGarbage() {
	garbageLock.Lock()
	defer garbageLock.Unlock()
	garbage = nil
}
