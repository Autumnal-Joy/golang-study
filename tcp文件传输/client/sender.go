package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// 获取命令行参数
	args := os.Args
	if len(args) != 2 {
		fmt.Println("go run sender.go [filename]")
		return
	}
	// 获取文件路径
	filePath := args[1]

	// 获取文件详细信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("os.Stat err: ", err)
		return
	}

	// 发起连接
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return
	}
	defer conn.Close()
	fmt.Println("服务器连接成功")

	// 发送文件名
	_, err = conn.Write([]byte(fileInfo.Name()))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return
	}
	buf := make([]byte, 4096)

	// 接收响应
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("os.Stat err: ", err)
		return
	}

	// 校对响应发送文件
	if string(buf[:n]) == "filename received" {
		sendFile(conn, filePath)
	}
}

func sendFile(conn net.Conn, filePath string) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("os.Open err: ", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 4096)
	for {
		// 读取本地文件
		n, err := file.Read(buf)
		if err == io.EOF {
			fmt.Println("文件发送完毕")
			return
		} else if err != nil {
			fmt.Println("file.Read err: ", err)
			return
		}

		// 发送文件
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("os.Write err: ", err)
			return
		}
	}
}
