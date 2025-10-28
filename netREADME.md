# 网络编程
## TCP协议
TCP协议是面向连接的，可靠的，基于字节流的协议。

服务端代码编写：
- 监听端口：lis, err := net.Listen("tcp", ":8080")
- 接收请求：conn, err := lis.Accept()
- 处理连接：go process(conn)