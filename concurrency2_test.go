package main

import (
	"sync"
	"testing"
)

// 有四个协程，分别打印A、B、C、D，要求按顺序打印ABCDABCD...，总共打印10次

func TestConcurrency2(t *testing.T) {
	var wg sync.WaitGroup
	channels := make([]chan byte, 4)
	over := make(chan struct{})
	final := make(chan struct{})
	for i := range channels {
		channels[i] = make(chan byte)
	}
	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go func(n int, ch <-chan byte) {
			defer wg.Done()
			for {
				select {
				case char := <-ch:
					t.Logf("goroutine %d: -> %c", n, char)
					over <- struct{}{}
				case <-final:
					return
				}
			}
		}(i, channels[i-1])
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 4; j++ {
			channels[j] <- (65 + byte(j))
			<-over
		}
	}
	wg.Wait()
	close(final)
	for i := range channels {
		close(channels[i])
	}
	close(over)
}
