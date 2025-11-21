package main

import (
	"context"
	"errors"
	"time"
)

// 设计一个带有超时功能的对象池
type ObjectPool[T any] struct {
	pool    chan T
	timeout time.Duration
	newFn   func() T
}

func NewObjectPool[T any](size int, timeout time.Duration, fn func() T) *ObjectPool[T] {
	objectPool := &ObjectPool[T]{
		pool:    make(chan T, size),
		timeout: timeout,
		newFn:   fn,
	}
	for range size {
		objectPool.pool <- fn()
	}
	return objectPool
}

func (op *ObjectPool[T]) Get(ctx context.Context) (T, error) {
	select {
	case obj := <-op.pool:
		return obj, nil
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case <-time.After(op.timeout):
		var zero T
		return zero, errors.New("Get Value timeout")
	}
}

func (op *ObjectPool[T]) Put(obj T) {
	select {
	case op.pool <- obj:
	default:
	}
}
