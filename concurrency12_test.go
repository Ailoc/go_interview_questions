package main

import (
	"sync"
	"testing"
)

// 给定两个线程，交替打印数字和字母，直到打印完1-26和A-Z。
func TestConcurrency12(t *testing.T) {
	numChan := make(chan struct{})
	charChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 26; i++ {
			<-numChan
			t.Logf("%d\n", i)
			charChan <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 'A'; i <= 'Z'; i++ {
			<-charChan
			t.Logf("%c\n", i)
			if i == 'Z' {
				break
			}
			numChan <- struct{}{}
		}
	}()
	numChan <- struct{}{}
	wg.Wait()
	close(numChan)
	close(charChan)
}
