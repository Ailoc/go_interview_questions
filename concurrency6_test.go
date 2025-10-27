package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

// 设计一个限流器，同一时间最多允许5个并发请求的测试用例

type RateLimiter6 struct {
	ch chan struct{}
}

func NewRateLimiter6(maxConcurrent int) *RateLimiter6 {
	rl := &RateLimiter6{
		ch: make(chan struct{}, maxConcurrent),
	}
	for i := 0; i < maxConcurrent; i++ {
		rl.ch <- struct{}{}
	}
	return rl
}
func (rl *RateLimiter6) Stop() {
	close(rl.ch)
}
func (rl *RateLimiter6) Acquire(ctx context.Context) bool { // 引入超时机制，防止大量请求阻塞
	select {
	case <-rl.ch:
		return true
	case <-ctx.Done():
		return false
	}
}
func (rl *RateLimiter6) Release() {
	rl.ch <- struct{}{}
}

func TestConcurrency6(t *testing.T) {
	rl := NewRateLimiter6(5)
	defer rl.Stop()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()
		go func(id int, rl *RateLimiter6, ctx context.Context) {
			defer wg.Done()
			if rl.Acquire(ctx) {
				t.Logf("Request %d is being processed", id)
				time.Sleep(20 * time.Millisecond)
				rl.Release()
			} else {
				t.Logf("Request %d timed out", id)
			}
		}(i, rl, ctx)
	}
	wg.Wait()
}
