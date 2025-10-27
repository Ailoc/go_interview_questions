package main

import (
	"sync"
	"testing"
	"time"
)

// 设计一个并发安全的限流器，限制每秒最多接收3个请求的测试用例

type RateLimiter struct {
	limit   int
	tickets chan struct{}
	ticker  *time.Ticker
	stop    chan struct{}
}

func NewRateLimiter(t time.Duration, limit int) *RateLimiter {
	rl := &RateLimiter{
		limit:   limit,
		tickets: make(chan struct{}, limit),
		ticker:  time.NewTicker(time.Second * t / time.Duration(limit)),
		stop:    make(chan struct{}),
	}
	for i := 0; i < limit; i++ {
		rl.tickets <- struct{}{}
	}
	go rl.refill()
	return rl
}

func (rl *RateLimiter) refill() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case rl.tickets <- struct{}{}:
			default:
			}
		case <-rl.stop:
			return
		}
	}
}
func (rl *RateLimiter) Acquire() {
	<-rl.tickets
}
func (rl *RateLimiter) Release() {
	select {
	case rl.tickets <- struct{}{}:
	default:
	}
}
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
	close(rl.stop)
}

func TestConcurrency5(t *testing.T) {
	rl := NewRateLimiter(1, 3)
	defer rl.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rl.Acquire()
			t.Logf("Request %d processed at %v", id, time.Now())
			time.Sleep(100 * time.Millisecond)
			// rl.Release()        // 及时归还令牌，允许限流器应对突发流量，此时每秒可处理超过3个请求
		}(i)
	}
	wg.Wait()
}
