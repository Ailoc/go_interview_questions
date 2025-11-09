package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// 实现一个固定数量的工作线程池，任务提交不阻塞，任务满时返回错误，支持优雅关闭
// pool := NewPool(5, 100)
// pool.Submit(task)
// pool.Shutdown()
type Pool struct {
	tasks  chan func()
	closed uint32
	wg     sync.WaitGroup
}

func NewPool(workers, queueSize int) *Pool {
	pool := &Pool{
		tasks: make(chan func(), queueSize),
	}
	pool.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer pool.wg.Done()

			for task := range pool.tasks {
				task()
			}
		}()
	}
	return pool
}

func (pool *Pool) Submit(task func()) error {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return fmt.Errorf("task queue is closed.")
	}
	select {
	case pool.tasks <- task:
		return nil
	default:
		return fmt.Errorf("task queue is full.")
	}
}
func (pool *Pool) Shutdown() {
	atomic.StoreUint32(&pool.closed, 1)
	close(pool.tasks)
	pool.wg.Wait()
}
func TestConcurrency18(t *testing.T) {
	pool := NewPool(5, 10)
	for i := 0; i < 20; i++ {
		err := pool.Submit(func() {
			fmt.Printf("task %d done\n", i)
		})
		if err != nil {
			t.Error(err)
		}
	}
	pool.Shutdown()
}
