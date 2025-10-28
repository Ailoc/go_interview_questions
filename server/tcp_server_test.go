package main

import (
	"bufio"
	"log"
	"net"
	"testing"
)

// go test -run TestTCPServer ./server -v
func TestTCPServer(t *testing.T) {
	// 1. 监听端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("监听端口出错: %v", err)
	}
	defer listener.Close()

	for {
		// 2. 接受连接
		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("接受连接出错: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			conn.Close()
			return
		}
		log.Printf("接收到数据：%s", string(buf[:n]))
		// 写回数据
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("写回数据出错: %v", err)
			return
		}
	}
}
