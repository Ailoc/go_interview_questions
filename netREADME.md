# 网络编程
## TCP协议
TCP协议是面向连接的，可靠的，基于字节流的协议。

服务端代码编写：
- 监听端口：

        lis, err := net.Listen("tcp", ":8080")
- 接收请求：

        conn, err := lis.Accept()
- 处理连接：

        go process(conn)

客户端代码：
- 建立连接：

        conn, err := net.Dial("tcp", "127.0.0.1")

- 发送数据：

            writer := bufio.NewWriter(conn)
            writer.Write([]byte{"data"})
            writer.Flush()
## UDP
UDP是一种无连接的协议，就是说在传输数据之前不需要建立连接

服务端代码编写：
- 监听端口：

        conn, err := net.ListenUDP()
- 接收数据：

        n, ip, err := conn.ReadFromUDP([]byte)
- 发送数据：

        _, err := conn.WriteToUDP([]byte, ip)
