package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 实现一个生产者消费者系统，三个生产者随机向channel发送整数，
// 两个消费者从channel中取出整数并打印总和，缓冲channel大小为5，生产1微秒，
// 如果是生产指定数量的数字，则需要外部控制，因为消费者在持续消费
func TestConcurrency(t *testing.T) {
	storageNum := make(chan int, 10)
	var mu sync.Mutex
	var wg sync.WaitGroup
	stop := make(chan struct{})
	Sum := 0
	wg.Add(5)
	for i := 0; i < 3; i++ {
		go Producer(storageNum, &wg, stop)
	}
	for i := 0; i < 2; i++ {
		go Consumer(storageNum, &Sum, &mu, &wg)
	}
	time.Sleep(time.Microsecond)
	close(stop)
	close(storageNum)
	wg.Wait()
	t.Logf("最终总和为：%d", Sum)

}
func Producer(ch chan int, wg *sync.WaitGroup, stop chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-stop:
			return
		default:
			randNum := rand.Intn(10)
			ch <- randNum
			fmt.Printf("生产数字: %d\n", randNum)
		}
	}
}
func Consumer(ch chan int, sum *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	localSum := 0
	for num := range ch {
		localSum += num
	}
	mu.Lock()
	*sum += localSum
	mu.Unlock()
}
