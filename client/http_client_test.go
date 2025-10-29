package main

import (
	"io"
	"net/http"
	"testing"
)

func TestHTTPClient(t *testing.T) {
	response, err := http.Get("http://127.0.0.1:8080")
	if err != nil {
		t.Fatalf("发送HTTP请求出错: %v", err)
	}
	defer response.Body.Close()

	t.Logf("响应状态码: %d", response.StatusCode)
	t.Logf("响应Content-Type: %s", response.Header.Get("Content-Type"))
	body := make([]byte, 1024)
	for {
		_, err := response.Body.Read(body)
		if err != nil {
			if err == io.EOF {
				t.Logf("读取响应体结束")
				break
			}
			t.Errorf("读取响应体出错: %v", err)
		}
	}

	t.Logf("响应体为: %s", string(body))
}
