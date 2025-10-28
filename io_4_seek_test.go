package main

import (
	"os"
	"testing"
)

func TestSeeker(t *testing.T) {
	file, err := os.OpenFile("./io_test.txt", os.O_RDWR, 0655) // 以读写模式打开文件,0655分别表示特殊权限，所有者权限，组权限，其他用户权限
	if err != nil {
		t.Fatalf("打开文件出错: %v", err)
	}
	defer file.Close()

	// 指定偏移量
	file.Seek(0, 2) // 从文件末尾开始偏移0个字节，即定位到文件末尾
	data := []byte("\n在文件末尾追加一行内容，用于测试Seek方法。")
	n, err := file.Write(data)
	if err != nil {
		t.Errorf("写入文件出错: %v", err)
	}
	t.Logf("成功写入 %d 字节到文件末尾", n)

	// 读取指定位置的数据
	file.Seek(0, 0) // 从文件开头开始偏移0个字节，即定位到文件开头
	buf := make([]byte, 64)
	n, err = file.Read(buf) // 实际读取了一部分数据
	if err != nil {
		t.Errorf("读取文件出错: %v", err)
	}
	t.Logf("从文件开头读取到的数据: %s", string(buf[:n]))
}
