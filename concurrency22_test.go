package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 使用channel实现一个定时器,类似time.After
func After(d time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(d)
		close(ch)
	}()
	return ch
}

func TestConcurrency22(t *testing.T) {
	t.Run("基本功能测试", func(t *testing.T) {
		delay := 100 * time.Millisecond
		start := time.Now()

		ch := After(delay)
		<-ch // 等待 channel 关闭

		elapsed := time.Since(start)

		// 允许 50ms 的误差
		if elapsed < delay || elapsed > delay+50*time.Millisecond {
			t.Errorf("期望延迟约 %v，实际延迟 %v", delay, elapsed)
		}

		t.Logf("基本功能测试通过，延迟 %v", elapsed)
	})

	t.Run("零延迟测试", func(t *testing.T) {
		start := time.Now()

		ch := After(0)
		<-ch

		elapsed := time.Since(start)

		// 零延迟应该几乎立即返回
		if elapsed > 50*time.Millisecond {
			t.Errorf("零延迟测试失败，实际延迟 %v", elapsed)
		}

		t.Logf("零延迟测试通过，延迟 %v", elapsed)
	})

	t.Run("极短延迟测试", func(t *testing.T) {
		delay := 1 * time.Millisecond
		start := time.Now()

		ch := After(delay)
		<-ch

		elapsed := time.Since(start)

		if elapsed > 100*time.Millisecond {
			t.Errorf("极短延迟测试失败，实际延迟 %v", elapsed)
		}

		t.Logf("极短延迟测试通过，延迟 %v", elapsed)
	})

	t.Run("较长延迟测试", func(t *testing.T) {
		delay := 500 * time.Millisecond
		start := time.Now()

		ch := After(delay)
		<-ch

		elapsed := time.Since(start)

		if elapsed < delay || elapsed > delay+100*time.Millisecond {
			t.Errorf("期望延迟约 %v，实际延迟 %v", delay, elapsed)
		}

		t.Logf("较长延迟测试通过，延迟 %v", elapsed)
	})

	t.Run("channel 关闭验证", func(t *testing.T) {
		ch := After(50 * time.Millisecond)

		// 等待 channel 关闭
		<-ch

		// 验证 channel 已关闭 - 从已关闭的 channel 读取会立即返回零值
		select {
		case _, ok := <-ch:
			if ok {
				t.Error("channel 应该已经关闭")
			}
		default:
			t.Error("从已关闭的 channel 读取应该立即返回")
		}

		t.Log("channel 关闭验证通过")
	})

	t.Run("多次读取已关闭的 channel", func(t *testing.T) {
		ch := After(50 * time.Millisecond)
		<-ch // 等待关闭

		// 多次读取已关闭的 channel 应该都能立即返回
		for i := 0; i < 5; i++ {
			select {
			case <-ch:
				// 正常，已关闭的 channel 可以读取多次
			case <-time.After(10 * time.Millisecond):
				t.Error("从已关闭的 channel 读取应该立即返回")
			}
		}

		t.Log("多次读取已关闭的 channel 测试通过")
	})

	t.Run("在 select 中使用", func(t *testing.T) {
		delay := 100 * time.Millisecond
		start := time.Now()

		select {
		case <-After(delay):
			elapsed := time.Since(start)
			if elapsed < delay || elapsed > delay+50*time.Millisecond {
				t.Errorf("期望延迟约 %v，实际延迟 %v", delay, elapsed)
			}
			t.Logf("select 语句测试通过，延迟 %v", elapsed)
		case <-time.After(200 * time.Millisecond):
			t.Error("After 应该在 200ms 内触发")
		}
	})

	t.Run("多个定时器竞速", func(t *testing.T) {
		ch1 := After(50 * time.Millisecond)
		ch2 := After(100 * time.Millisecond)
		ch3 := After(150 * time.Millisecond)

		start := time.Now()
		var first, second, third time.Duration

		// 应该按顺序触发
		<-ch1
		first = time.Since(start)

		<-ch2
		second = time.Since(start)

		<-ch3
		third = time.Since(start)

		if first > 100*time.Millisecond {
			t.Errorf("第一个定时器触发过晚: %v", first)
		}
		if second < 100*time.Millisecond || second > 150*time.Millisecond {
			t.Errorf("第二个定时器触发时间异常: %v", second)
		}
		if third < 150*time.Millisecond || third > 200*time.Millisecond {
			t.Errorf("第三个定时器触发时间异常: %v", third)
		}

		t.Logf("多个定时器竞速测试通过: %v, %v, %v", first, second, third)
	})

	t.Run("并发创建多个定时器", func(t *testing.T) {
		count := 50
		delay := 100 * time.Millisecond

		var wg sync.WaitGroup
		start := time.Now()

		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				ch := After(delay)
				<-ch
			}(i)
		}

		wg.Wait()
		elapsed := time.Since(start)

		// 所有定时器应该几乎同时触发
		if elapsed < delay || elapsed > delay+100*time.Millisecond {
			t.Errorf("期望所有定时器在约 %v 后触发，实际耗时 %v", delay, elapsed)
		}

		t.Logf("并发创建 %d 个定时器测试通过，耗时 %v", count, elapsed)
	})

	t.Run("不同延迟的并发定时器", func(t *testing.T) {
		delays := []time.Duration{
			50 * time.Millisecond,
			100 * time.Millisecond,
			150 * time.Millisecond,
			200 * time.Millisecond,
			250 * time.Millisecond,
		}

		var wg sync.WaitGroup
		results := make([]time.Duration, len(delays))

		start := time.Now()

		for i, delay := range delays {
			wg.Add(1)
			go func(idx int, d time.Duration) {
				defer wg.Done()
				ch := After(d)
				<-ch
				results[idx] = time.Since(start)
			}(i, delay)
		}

		wg.Wait()

		// 验证每个定时器的触发时间
		for i, delay := range delays {
			elapsed := results[i]
			if elapsed < delay || elapsed > delay+100*time.Millisecond {
				t.Errorf("定时器 %d (延迟 %v) 触发时间异常: %v", i, delay, elapsed)
			}
		}

		t.Log("不同延迟的并发定时器测试通过")
	})

	t.Run("select 中的超时模式", func(t *testing.T) {
		// 模拟一个永远不会发送数据的 channel
		slowCh := make(chan int)

		timeout := 100 * time.Millisecond
		start := time.Now()

		select {
		case <-slowCh:
			t.Error("不应该从 slowCh 接收到数据")
		case <-After(timeout):
			elapsed := time.Since(start)
			if elapsed < timeout || elapsed > timeout+50*time.Millisecond {
				t.Errorf("超时时间不准确，期望 %v，实际 %v", timeout, elapsed)
			}
			t.Logf("超时模式测试通过，超时时间 %v", elapsed)
		}
	})

	t.Run("多个超时竞争", func(t *testing.T) {
		start := time.Now()
		var winner string

		select {
		case <-After(50 * time.Millisecond):
			winner = "50ms"
		case <-After(100 * time.Millisecond):
			winner = "100ms"
		case <-After(150 * time.Millisecond):
			winner = "150ms"
		}

		elapsed := time.Since(start)

		if winner != "50ms" {
			t.Errorf("应该是 50ms 的定时器先触发，实际是 %s", winner)
		}

		if elapsed > 100*time.Millisecond {
			t.Errorf("最短的定时器应该在 100ms 内触发，实际 %v", elapsed)
		}

		t.Logf("多个超时竞争测试通过，winner: %s, 耗时 %v", winner, elapsed)
	})

	t.Run("高并发压力测试", func(t *testing.T) {
		count := 1000
		delay := 10 * time.Millisecond

		var wg sync.WaitGroup
		var completed int32

		start := time.Now()

		for i := 0; i < count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ch := After(delay)
				<-ch
				atomic.AddInt32(&completed, 1)
			}()
		}

		wg.Wait()
		elapsed := time.Since(start)

		if completed != int32(count) {
			t.Errorf("期望完成 %d 个定时器，实际完成 %d 个", count, completed)
		}

		t.Logf("高并发压力测试通过，创建并等待 %d 个定时器，耗时 %v", count, elapsed)
	})

	t.Run("定时器精度测试", func(t *testing.T) {
		delays := []time.Duration{
			10 * time.Millisecond,
			50 * time.Millisecond,
			100 * time.Millisecond,
			200 * time.Millisecond,
		}

		for _, delay := range delays {
			start := time.Now()
			<-After(delay)
			elapsed := time.Since(start)

			// 计算误差百分比
			diff := elapsed - delay
			if diff < 0 {
				diff = -diff
			}

			// 允许 20% 的误差或最多 50ms
			maxError := delay / 5
			if maxError > 50*time.Millisecond {
				maxError = 50 * time.Millisecond
			}

			if diff > maxError {
				t.Errorf("延迟 %v: 精度不足，实际 %v，误差 %v", delay, elapsed, diff)
			}
		}

		t.Log("定时器精度测试通过")
	})

	t.Run("返回只读 channel 验证", func(t *testing.T) {
		ch := After(10 * time.Millisecond)

		// 验证类型是只读 channel
		var _ <-chan struct{} = ch

		// 这行代码如果取消注释会编译错误，因为 ch 是只读的
		// ch <- struct{}{} // 编译错误: cannot send to receive-only channel

		<-ch
		t.Log("返回只读 channel 验证通过")
	})

	t.Run("goroutine 泄漏检查", func(t *testing.T) {
		// 创建定时器但不等待
		for i := 0; i < 10; i++ {
			After(1 * time.Millisecond)
		}

		// 等待足够长的时间让所有定时器触发
		time.Sleep(50 * time.Millisecond)

		// 注意：这个测试只是基本检查，真正的 goroutine 泄漏检查需要工具
		t.Log("goroutine 泄漏检查完成（基础检查）")
	})

	t.Run("连续创建和等待", func(t *testing.T) {
		iterations := 20
		delay := 10 * time.Millisecond

		start := time.Now()

		for i := 0; i < iterations; i++ {
			ch := After(delay)
			<-ch
		}

		elapsed := time.Since(start)

		// 总时间应该约等于 iterations * delay
		expectedMin := time.Duration(iterations) * delay
		expectedMax := expectedMin + 200*time.Millisecond

		if elapsed < expectedMin || elapsed > expectedMax {
			t.Errorf("连续创建和等待时间异常，期望 %v-%v，实际 %v",
				expectedMin, expectedMax, elapsed)
		}

		t.Logf("连续创建和等待测试通过，%d 次迭代耗时 %v", iterations, elapsed)
	})

	t.Run("与 time.After 对比", func(t *testing.T) {
		delay := 100 * time.Millisecond

		// 测试我们的实现
		start1 := time.Now()
		<-After(delay)
		elapsed1 := time.Since(start1)

		// 测试标准库
		start2 := time.Now()
		<-time.After(delay)
		elapsed2 := time.Since(start2)

		// 两者时间应该接近
		diff := elapsed1 - elapsed2
		if diff < 0 {
			diff = -diff
		}

		if diff > 50*time.Millisecond {
			t.Logf("注意：与标准库 time.After 有较大差异，自定义: %v, 标准: %v", elapsed1, elapsed2)
		} else {
			t.Logf("与 time.After 对比通过，自定义: %v, 标准: %v", elapsed1, elapsed2)
		}
	})
}
