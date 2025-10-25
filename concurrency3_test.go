package main

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

// 实现一个worker pool，限制最多3个goroutine同时工作，处理10个任务，每个任务计算任务的平方并打印结果
func TestConcurrency3(t *testing.T) {
	task := make(chan int, 10)
	for i := 0; i < 10; i++ {
		task <- i
	}
	close(task)
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go worker(i, task, &wg)
	}
	wg.Wait()
	t.Log("all tasks done")
}
func worker(id int, task <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for num := range task {
		result := math.Pow(float64(num), 2)
		fmt.Printf("goroutine %d: %d square result %f\n", id, num, result)
	}
}
