package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// 实现一个Once函数，确保传入的函数只会被执行一次，
type Once struct {
	mu   sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}
func (o *Once) doSlow(f func()) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
func TestConcurrency15(t *testing.T) {

}
