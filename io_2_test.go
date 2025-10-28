package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestIO2(t *testing.T) {
	file, err := os.Open("io_test.txt") // 默认为只读模式
	if err != nil {
		t.Fatalf("打开文件出错: %v", err)
	}
	defer file.Close()

	buf := make([]byte, 16) // 创建一个字节数组，长度为16
	all := make([]byte, 0)
	for {
		n, err := file.Read(buf) // 读取数据到buf中，返回读取的字节数和错误信息
		if err != nil {
			if err == io.EOF {
				t.Logf("读取到文件末尾，结束读取")
				fmt.Printf("读到的数据为：%s", string(all))
				break
			}
			t.Errorf("读取数据出错: %v", err)
			break
		}
		all = append(all, buf[:n]...)
	}
}
