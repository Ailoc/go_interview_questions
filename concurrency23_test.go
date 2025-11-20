package main

import (
	"io"
	"net/http"
	"sync"
)

// 并发爬取1000个url，限制最大并发数为50个，将所有结果汇总
func crawl(urls []string) ([]string, error) {
	sem := make(chan struct{}, 50)
	tasks := make(chan string)
	results := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range tasks {
				sem <- struct{}{}
				res, _ := http.Get(url)
				body, _ := io.ReadAll(res.Body)
				results <- string(body)
				<-sem
			}
		}()
	}

	go func() {
		for _, url := range urls {
			tasks <- url
		}
		close(tasks)
	}()

	wg.Wait()
	close(results)
	var res []string
	for r := range results {
		res = append(res, r)
	}
	return res, nil
}
