package main

import (
	"bufio"
	"net"
	"testing"
)

func TestUDPClient(t *testing.T) {
	// 1. 连接服务器
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9090,
	})
	if err != nil {
		t.Fatalf("连接服务器出错: %v", err)
	}
	defer conn.Close()
	data := []byte("Hello UDP Server")
	// 2. 发送固定测试数据
	writer := bufio.NewWriter(conn)
	n, err := writer.Write(data)
	if err != nil {
		t.Errorf("发送数据出错: %v", err)
	}
	t.Logf("成功发送 %d 字节数据到缓冲区", n)
	err = writer.Flush()
	if err != nil {
		t.Errorf("刷新数据出错: %v", err)
	}
	// 3. 接收服务器响应
	var buf [1024]byte
	reader := bufio.NewReader(conn)
	n, err = reader.Read(buf[:])
	if err != nil {
		t.Errorf("接收服务器响应出错: %v", err)
	}
	t.Logf("接收到服务器响应：%s", string(buf[:n]))
}
