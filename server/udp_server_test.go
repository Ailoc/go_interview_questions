package main

import (
	"net"
	"testing"
)

func TestUDPServer(t *testing.T) {
	// 1. 监听端口
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9090,
	})
	if err != nil {
		t.Fatalf("监听端口出错: %v", err)
	}
	defer conn.Close()

	for {
		var buf [1024]byte
		n, ip, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			t.Errorf("读取数据出错: %v", err)
			continue
		}
		t.Logf("接收到来自 %s 的数据：%s", ip.String(), string(buf[:n]))
		// 写回数据
		_, err = conn.WriteToUDP(buf[:n], ip)
		if err != nil {
			t.Errorf("写回数据出错: %v", err)
			continue
		}
		t.Logf("成功写回 %d 字节数据到 %s", n, ip.String())
	}
}
