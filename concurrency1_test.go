package main

import (
	"sync"
	"testing"
)

// 代码的要求是开启100个协程顺序打印1-1000，且保证协程号为n的，打印尾数为n的数字
func TestConcurrency1(t *testing.T) {
	channels := make([]chan int, 100)
	over := make(chan struct{})
	final := make(chan struct{})
	var wg sync.WaitGroup
	for i := range channels {
		channels[i] = make(chan int)
	}
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(n int, ch <-chan int) {
			defer wg.Done()
			for {
				select {
				case num := <-ch:
					t.Logf("Goroutine %d received number %d", n, num)
					over <- struct{}{}
				case <-final:
					return
				}
			}
		}(i, channels[i-1])
	}
	g_num := 0
	for i := 1; i <= 1000; i++ {
		g_num = i % 100
		if g_num == 0 {
			g_num = 100
		}
		channels[g_num-1] <- i
		<-over
	}
	close(final)
	wg.Wait()
}
