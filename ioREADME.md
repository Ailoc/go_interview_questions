# IO
## io包主要包含以下接口
- io.Reader, 从数据源读取到字节切片
- io.Writer, 将数据写入目标
- io.Closer, 关闭文件
- io.Seeker, 定位数据源的位置,设置读取和写入时的偏移量

```go
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}
// whence: 起始值
const (
    SeekStart = 0
    SeekCurrent = 1
    SeekEnd = 2
)
```
### 方法
- os.ReadFile(name string) ([]byte, error)
- os.WriteFile(name string, data []byte, perm fs.FileMode)
- io.ReadAll(r io.Reader) ([]byte, error)  -> ReadAll(file) 
- ...

## 文件操作API
- Create(name string) (file *File, err Error)
- Open(name string) (file *File, err Error)
- OpenFile(name string, flag int, perm uint32) (file *File, err Error)
- Write(b []byte) (n int, err Error)
- WriteAt(b []byte, off int64) (n int, err Error)   -> 从指定位置读
- WriteString(s string) (ret int, err Error)
- Read(b []byte) (n int, err Error)
- ReadAt(b []byte, off int64) (n int, err Error)
- Remove(name string) Error   -> 删除文件

## bufio，自带缓冲区的io
- bufio.Reader  -> 包含reader.ReadLine()方法读取一行
- bufio.Writer