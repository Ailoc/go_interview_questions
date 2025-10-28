package main

import (
	"os"
	"testing"
)

func TestIO3(t *testing.T) {
	file, err := os.Create("./io_test_write.txt")
	if err != nil {
		t.Fatalf("创建文件出错: %v", err)
	}
	defer file.Close()

	data := []byte("这是要写入文件的内容，用于测试io写操作。\n第二行内容。\n第三行内容。")
	n, err := file.Write(data)
	if err != nil {
		t.Errorf("写入文件出错: %v", err)
	}
	t.Logf("成功写入 %d 字节到文件", n)
}
