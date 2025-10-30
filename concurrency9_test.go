package main

import (
	"sync"
	"testing"
)

// 开启两个协程，交替打印奇数和偶数，从1打印到20。
func TestConcurrency9(t *testing.T) {
	oddChan := make(chan struct{})
	evenChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 20; i += 2 {
			<-oddChan
			t.Logf("奇数：%d\n", i)
			evenChan <- struct{}{}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 2; i <= 20; i += 2 {
			<-evenChan
			t.Logf("偶数：%d\n", i)
			if i == 20 {
				break
			}
			oddChan <- struct{}{}
		}
	}()
	oddChan <- struct{}{}
	wg.Wait()
	close(oddChan)
	close(evenChan)
}
