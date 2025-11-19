package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 实现一个带超时功能的WorkPool，提交任务超时时放弃提交
//
//

type WorkPool struct {
	tasks        chan func() error
	workers      uint32
	shutdown     chan struct{}
	shutdownOnce *sync.Once
	closed       atomic.Bool
	taskOnce     *sync.Once
	wg           *sync.WaitGroup
}

func NewWorkPool(workerNum, taskNum uint32) *WorkPool {
	if workerNum <= 0 {
		workerNum = 5
	}
	if taskNum <= 0 {
		taskNum = 0
	}
	pool := &WorkPool{
		tasks:        make(chan func() error, taskNum),
		workers:      workerNum,
		shutdown:     make(chan struct{}),
		shutdownOnce: new(sync.Once),
		taskOnce:     new(sync.Once),
		wg:           new(sync.WaitGroup),
	}
	pool.Start()
	return pool
}

func (pool *WorkPool) Start() {
	for i := 0; i < int(pool.workers); i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()

			for task := range pool.tasks {
				if task != nil {
					_ = task()
				}
			}
		}()
	}
}

func (pool *WorkPool) Close() {
	pool.shutdownOnce.Do(func() {
		pool.closed.Store(true)
		close(pool.shutdown)
	})
}

func (pool *WorkPool) Submit(task func() error, timeout time.Duration) bool {
	if pool.closed.Load() {
		return false
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return false
	case <-pool.shutdown:
		return false
	case pool.tasks <- task:
		return true
	}
}

func (pool *WorkPool) Wait() {
	pool.Close()
	pool.taskOnce.Do(func() {
		close(pool.tasks)
	})
	pool.wg.Wait()
}

func TestConcurrency20(t *testing.T) {
	t.Run("基本功能测试", func(t *testing.T) {
		pool := NewWorkPool(3, 10)
		defer pool.Wait()

		var counter int32
		var wg sync.WaitGroup

		// 提交 10 个任务
		for i := 0; i < 10; i++ {
			wg.Add(1)
			task := func() error {
				atomic.AddInt32(&counter, 1)
				time.Sleep(10 * time.Millisecond)
				wg.Done()
				return nil
			}
			if !pool.Submit(task, time.Second) {
				t.Errorf("任务提交失败")
				wg.Done()
			}
		}

		wg.Wait()
		if counter != 10 {
			t.Errorf("期望执行 10 个任务，实际执行了 %d 个", counter)
		}
		t.Logf("成功执行了 %d 个任务", counter)
	})

	t.Run("超时功能测试", func(t *testing.T) {
		// 创建一个队列大小为 0 的 pool，这样如果没有 worker 立即处理，就会超时
		pool := NewWorkPool(1, 0)

		// 先让 worker 忙碌
		var blockChan = make(chan struct{})
		blockingTask := func() error {
			<-blockChan // 阻塞住唯一的 worker
			return nil
		}
		pool.Submit(blockingTask, time.Second)

		// 尝试提交第二个任务，应该超时
		fastTimeout := func() error {
			t.Log("这个任务不应该被执行")
			return nil
		}

		submitted := pool.Submit(fastTimeout, 10*time.Millisecond)
		if submitted {
			t.Errorf("期望任务提交超时，但实际成功了")
		} else {
			t.Log("任务提交超时，符合预期")
		}

		close(blockChan) // 释放阻塞
		pool.Wait()
	})

	t.Run("关闭后拒绝新任务", func(t *testing.T) {
		pool := NewWorkPool(3, 10)

		var counter int32
		task := func() error {
			atomic.AddInt32(&counter, 1)
			return nil
		}

		// 提交一些任务
		for i := 0; i < 5; i++ {
			pool.Submit(task, time.Second)
		}

		// 关闭 pool
		pool.Close()

		// 尝试提交新任务，应该失败
		submitted := pool.Submit(task, time.Second)
		if submitted {
			t.Errorf("期望关闭后拒绝新任务，但实际接受了")
		} else {
			t.Log("关闭后成功拒绝新任务")
		}

		pool.Wait()
		t.Logf("关闭前提交的任务数量：%d", counter)
	})

	t.Run("并发提交测试", func(t *testing.T) {
		pool := NewWorkPool(5, 100)
		defer pool.Wait()

		var counter int32
		var submitWg sync.WaitGroup
		concurrentSubmitters := 10
		tasksPerSubmitter := 10

		// 10 个 goroutine 同时提交任务
		for i := 0; i < concurrentSubmitters; i++ {
			submitWg.Add(1)
			go func(id int) {
				defer submitWg.Done()
				for j := 0; j < tasksPerSubmitter; j++ {
					task := func() error {
						atomic.AddInt32(&counter, 1)
						time.Sleep(time.Millisecond)
						return nil
					}
					if !pool.Submit(task, time.Second) {
						t.Errorf("Submitter %d: 任务 %d 提交失败", id, j)
					}
				}
			}(i)
		}

		submitWg.Wait()
		pool.Close()
		pool.Wait()

		expected := int32(concurrentSubmitters * tasksPerSubmitter)
		if counter != expected {
			t.Errorf("期望执行 %d 个任务，实际执行了 %d 个", expected, counter)
		} else {
			t.Logf("并发提交成功：%d 个任务全部执行完成", counter)
		}
	})

	t.Run("并发提交和关闭测试", func(t *testing.T) {
		pool := NewWorkPool(5, 50)

		var counter int32
		var submitWg sync.WaitGroup

		// 启动多个 goroutine 提交任务
		for i := 0; i < 5; i++ {
			submitWg.Add(1)
			go func() {
				defer submitWg.Done()
				for j := 0; j < 20; j++ {
					task := func() error {
						atomic.AddInt32(&counter, 1)
						time.Sleep(time.Millisecond)
						return nil
					}
					pool.Submit(task, 100*time.Millisecond)
					time.Sleep(time.Millisecond)
				}
			}()
		}

		// 同时有一个 goroutine 在稍后关闭 pool
		go func() {
			time.Sleep(10 * time.Millisecond)
			pool.Close()
		}()

		submitWg.Wait()
		pool.Wait()

		t.Logf("并发提交和关闭测试：成功执行了 %d 个任务", counter)
	})

	t.Run("Worker 数量限制测试", func(t *testing.T) {
		maxWorkers := 3
		pool := NewWorkPool(uint32(maxWorkers), 10)
		defer pool.Wait()

		var activeWorkers int32
		var maxActive int32
		var wg sync.WaitGroup

		// 提交比 worker 数量更多的任务
		for i := 0; i < 10; i++ {
			wg.Add(1)
			task := func() error {
				defer wg.Done()

				// 记录当前活跃的 worker 数量
				current := atomic.AddInt32(&activeWorkers, 1)

				// 更新最大活跃数量
				for {
					max := atomic.LoadInt32(&maxActive)
					if current <= max || atomic.CompareAndSwapInt32(&maxActive, max, current) {
						break
					}
				}

				// 模拟工作
				time.Sleep(50 * time.Millisecond)

				atomic.AddInt32(&activeWorkers, -1)
				return nil
			}
			pool.Submit(task, time.Second)
		}

		wg.Wait()

		if maxActive > int32(maxWorkers) {
			t.Errorf("期望最多 %d 个并发 worker，但观察到 %d 个", maxWorkers, maxActive)
		} else {
			t.Logf("Worker 数量限制正确：最大并发数 %d <= %d", maxActive, maxWorkers)
		}
	})

	t.Run("队列满时超时测试", func(t *testing.T) {
		workers := 1
		queueSize := 2
		pool := NewWorkPool(uint32(workers), uint32(queueSize))

		// 阻塞 worker
		blockChan := make(chan struct{})
		blockingTask := func() error {
			<-blockChan
			return nil
		}

		// 填满队列：1 个在执行，2 个在队列中
		pool.Submit(blockingTask, time.Second)
		pool.Submit(blockingTask, time.Second)
		pool.Submit(blockingTask, time.Second)

		// 现在队列满了，新任务应该超时
		submitted := pool.Submit(func() error {
			return nil
		}, 50*time.Millisecond)

		if submitted {
			t.Errorf("期望队列满时任务超时，但实际成功提交")
		} else {
			t.Log("队列满时任务超时，符合预期")
		}

		close(blockChan)
		pool.Wait()
	})

	t.Run("Wait 阻塞直到所有任务完成", func(t *testing.T) {
		pool := NewWorkPool(3, 10)

		var counter int32
		for i := 0; i < 10; i++ {
			task := func() error {
				time.Sleep(50 * time.Millisecond)
				atomic.AddInt32(&counter, 1)
				return nil
			}
			pool.Submit(task, time.Second)
		}

		start := time.Now()
		pool.Wait()
		duration := time.Since(start)

		if counter != 10 {
			t.Errorf("期望执行 10 个任务，实际执行了 %d 个", counter)
		}

		// 3 个 worker，每个任务 50ms，10 个任务至少需要 4 批次，大约 150ms
		if duration < 150*time.Millisecond {
			t.Errorf("Wait 过早返回，可能没有等待所有任务完成，耗时 %v", duration)
		} else {
			t.Logf("Wait 正确等待所有任务完成，耗时 %v", duration)
		}
	})
}
