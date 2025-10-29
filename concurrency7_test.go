package main

import (
	"sync"
	"testing"
)

// 实现Ping-Pong并发程序，使用无缓冲channel，交替执行10次打印"Ping"和"Pong"。
func TestConcurrency7(t *testing.T) {
	pingChan := make(chan struct{})
	pongChan := make(chan struct{})
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			<-pingChan
			t.Logf("-------------%d---------------", i+1)
			t.Log("Ping---->\n")
			pongChan <- struct{}{}
		}
		close(pongChan)
		close(stop)
	}()
	go func() {
		defer wg.Done()
		for range pongChan {
			// <-pongChan
			t.Log("<---Pong\n-------------------------------------")
			select {
			case pingChan <- struct{}{}:
			case <-stop:
				return
			}
		}
	}()
	// 启动Ping-Pong
	pingChan <- struct{}{}
	wg.Wait()
	close(pingChan)
	// close(pongChan)
}
