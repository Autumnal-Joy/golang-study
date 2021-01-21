package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	server()
}

func server() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("等待客户端连接")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			return
		}
		go handleConnect(conn)
	}
}

func handleConnect(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr()
	fmt.Println("连接客户端: ", addr)

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			fmt.Println("客户端关闭连接")
			return
		}
		if err != nil {
			fmt.Println("conn.Read err:", err)
			return
		}
		str := strings.TrimSpace(string(buf[:n]))
		if str == "exit" || str == "quit" || str == "q" {
			fmt.Printf("客户端结束通信\n")
			return
		}
		fmt.Printf("读取客户端数据: %q\n", str)

		res := strings.ToUpper(str)

		n, err = conn.Write([]byte(res))
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return
		}
	}
}
