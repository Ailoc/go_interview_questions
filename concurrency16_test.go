package main

import (
	"sync"
	"testing"
	"time"
)

// 实现一个限速的爬虫
func TestConcurrency16(t *testing.T) {
	rate := 5
	tokenChan := make(chan struct{}, rate)

	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	var wg sync.WaitGroup
	stopChan := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ticker.C {
			select {
			case tokenChan <- struct{}{}:
			case <-stopChan:
				return
			default:
			}
		}
	}()

	for i := 0; i < 10; i++ {
		<-tokenChan
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Logf("爬取页面 %d\n", i+1)
		}()
	}
	<-time.After(time.Second * 2)
	close(stopChan)
	wg.Wait()
}
