package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 在微服务中，需要并发调用多个下游任务（如用户服务，订单服务），任意一个失败就取消其他调用并返回错误
// 实现一个Group
// g := errgroup.WithContext(ctx)
// g.Go(func() error {...})
// err := g.Wait()  等待所有任务完成，返回第一个错误
type Group struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func WithContext(parentCtx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	return &Group{
		ctx:    ctx,
		cancel: cancel,
	}, ctx
}
func (g *Group) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		if err := g.ctx.Err(); err != nil {
			return
		}
		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

func TestConcurrency17(t *testing.T) {
	ctx := context.Background()
	g, ctx := WithContext(ctx)
	g.Go(func() error {
		time.Sleep(time.Second)
		return fmt.Errorf("user service failed.")
	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second):
			return nil
		}
	})
	if err := g.Wait(); err != nil {
		fmt.Println("Error:", err)
	}

}
