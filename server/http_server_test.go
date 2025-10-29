package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("请求方法为：%s", r.Method)
		fmt.Printf("请求Content-Type: %s", r.Header.Get("Content-Type"))
		fmt.Printf("请求URL为: %s", r.RequestURI)
		fmt.Printf("请求参数为: %s", r.URL.Query()) // r.URL.Query().Get("name")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("读取请求体出错: %v", err)
		}
		fmt.Printf("请求体为: %s", string(body))
		w.Write([]byte("Hello, HTTP Server!"))
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		t.Fatalf("启动HTTP服务器出错: %v", err)
	}
}
