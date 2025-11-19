package main

import (
	"sync"
	"testing"
	"time"
)

// 实现一个有界信号量
// 信号量是一种同步原语，用于控制对共享资源的访问。它维护一个计数器，表示可用资源的数量。
// 当一个goroutine需要访问资源时，它会尝试获取一个信号量。如果信号量的计数器大于零，则该goroutine可以获取信号量并继续执行。
// 如果信号量的计数器为零，则该goroutine会被阻塞，直到其他goroutine释放信号量。
// 信号量可以用于实现互斥锁、读写锁、生产者-消费者模型等并发控制机制。
//
// 有界信号量在信号量的基础上限制了信号量的最大值，以防止资源过度使用。
// Acquire(), Release()用于获取和释放信号量

type BoundedSemaphore struct {
	ch chan struct{}
}

func NewBoundedSemaphore(max int) *BoundedSemaphore {
	return &BoundedSemaphore{
		ch: make(chan struct{}, max),
	}
}
func (bs *BoundedSemaphore) Acquire() {
	bs.ch <- struct{}{}
}

func (bs *BoundedSemaphore) Release() {
	<-bs.ch
}

func TestConcurrency19(t *testing.T) {
	// 创建一个容量为3的有界信号量
	sem := NewBoundedSemaphore(3)

	// 用于记录同时运行的goroutine数量
	var running int
	var mu sync.Mutex
	var maxRunning int

	// 启动10个goroutine，但同时只能有3个运行
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem.Acquire()       // 获取信号量
			defer sem.Release() // 释放信号量

			mu.Lock()
			running++
			if running > maxRunning {
				maxRunning = running
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond) // 模拟工作

			mu.Lock()
			running--
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	// 验证最多只有3个goroutine同时运行
	if maxRunning != 3 {
		t.Errorf("期望最多3个goroutine同时运行，实际: %d", maxRunning)
	}
	t.Logf("测试通过，最大并发数: %d", maxRunning)

}
