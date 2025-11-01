package main

import (
	"sync"
	"testing"
)

// 并发计算数组中每个数的平方和
func TestConcurrency14(t *testing.T) {
	numArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := make([]int, len(numArray))
	chunkSize := (len(numArray) + 2) / 3 // 将数组分成3块
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(numArray) {
			end = len(numArray)
		}
		go func(arr, res *[]int, start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				(*res)[j] = (*arr)[j] * (*arr)[j]
			}
		}(&numArray, &result, start, end)
	}
	wg.Wait()
	t.Logf("平方结果: %v", result)
}
