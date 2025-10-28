package main

import (
	"net"
	"testing"
)

// go test -run TestTCPServer ./server -v
func TestTCPClient(t *testing.T) {
	// 1. 连接服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("连接服务器出错: %v", err)
	}
	defer conn.Close()

	// 2. 发送固定测试数据
	testMessages := []string{"Hello", "World", "Test Message"}

	for _, msg := range testMessages {
		t.Logf("发送: %s", msg)
		n, err := conn.Write([]byte(msg))
		if err != nil {
			t.Errorf("发送数据出错: %v", err)
			continue
		}
		t.Logf("成功发送 %d 字节数据到服务器", n)

		// 3. 接收服务器响应
		var buf [1024]byte
		n, err = conn.Read(buf[:])
		if err != nil {
			t.Errorf("接收服务器响应出错: %v", err)
			continue
		}
		t.Logf("接收到服务器响应：%s", string(buf[:n]))
	}
}
