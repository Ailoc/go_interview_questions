package main

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestIO(t *testing.T) {
	reader := strings.NewReader("这是一个字符串，准备读取到一个字节数组") // 返回一个reader结构体，实现了Reader接口
	buf := make([]byte, 12)                            // 创建一个字节数组，长度为12
	for {
		n, err := reader.Read(buf) // 读取数据到buf中，返回读取的字节数和错误信息
		if err != nil {
			if err == io.EOF { // 判断是否读取到文件末尾
				t.Logf("读取到文件末尾，结束读取")
				break
			}
			t.Errorf("读取数据出错: %v", err)
			break
		}
		fmt.Print(string(buf[:n])) // 将读取到的字节数组转换为字符串并打印
	}
}
