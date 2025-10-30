package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 给定n个任务，每个任务sleep任意时间并返回ID
// 使用最多M个并发worker执行，收集所有结果
func TestConcurrency10(t *testing.T) {
	numTasks := 20
	numWorkers := 5
	tasks := make(chan Task, numTasks)
	results := make(chan Result, numTasks)
	var wg sync.WaitGroup
	for i := 1; i <= numTasks; i++ {
		tasks <- Task{ID: i}
	}
	close(tasks)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go work(tasks, results, &wg)
	}
	wg.Wait()
	close(results) // ✅ 关闭results，这样for range才能退出
	for result := range results {
		t.Logf("任务 %d 完成，耗时 %v", result.ID, result.sleep)
	}
}

type Task struct {
	ID int
}
type Result struct {
	ID    int
	sleep time.Duration
}

func work(tasks chan Task, result chan Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		sleepTime := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(sleepTime)
		result <- Result{ID: task.ID, sleep: sleepTime}
	}
}
