package main

import (
	"strings"
	"sync"
	"testing"
)

// 给定一个大型日志文件，按级别分类统计数量，避免race condition

func TestConcurrency4(t *testing.T) {
	logs := []string{
		"INFO: User logged in",
		"ERROR: Database connection failed",
		"WARNING: Disk space low",
		"INFO: File uploaded",
		"ERROR: Timeout while fetching data",
		"INFO: User logged out",
		"WARNING: High memory usage detected",
		"ERROR: Failed to write to disk",
		"INFO: Scheduled job started",
		"WARNING: CPU temperature high",
	}
	var mu sync.Mutex
	counts := make(map[string]int)
	var wg sync.WaitGroup
	numWorkers := 3
	chunkSize := (len(logs) + numWorkers - 1) / numWorkers
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(logs) {
			end = len(logs)
		}
		wg.Add(1)
		go countLogLevels(&logs, start, end, &mu, &wg, &counts)
	}
	wg.Wait()
	for level, count := range counts {
		t.Logf("%s: %d", level, count)
	}

}
func countLogLevels(logs *[]string, start, end int, mu *sync.Mutex, wg *sync.WaitGroup, counts *map[string]int) {
	defer wg.Done()

	localCounts := make(map[string]int)
	for i := start; i < end; i++ {
		level := (*logs)[i][:strings.Index((*logs)[i], ":")]
		localCounts[level]++
	}

	mu.Lock()
	defer mu.Unlock()
	for level, count := range localCounts {
		(*counts)[level] += count
	}
}
