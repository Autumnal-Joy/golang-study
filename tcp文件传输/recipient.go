package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// 启动监听
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("等待连接")

	// 监听连接
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("listener.Accept err: ", err)
		return
	}
	defer conn.Close()

	// 接收文件名
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err: ", err)
		return
	}
	conn.Write([]byte("filename received"))

	receiveFile(conn, string(buf[:n]))
}

func receiveFile(conn net.Conn, filename string) {
	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		println("os.Create err: ", err)
		return
	}
	defer file.Close()

	// 保存文件到本地
	buf := make([]byte, 4096)
	for {
		// 接收数据
		n, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("文件接收完毕")
			return
		} else if err != nil {
			fmt.Println("conn.Read err: ", err)
			return
		}

		// 保存文件
		_, err = file.Write(buf[:n])
		if err != nil {
			println("file.Write err: ", err)
			return
		}
	}
}
