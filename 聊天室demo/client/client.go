package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	client()
}

func client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Dial err:", err)
	}
	defer conn.Close()

	buf := make([]byte, 4096)

	go func() {
		str := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(str)
			if err != nil {
				fmt.Println("os.Stdin.Read err: ", err)
				continue
			}
			conn.Write(str[:n])
		}
	}()

	for {
		n, err := conn.Read(buf)
		if n == 0 {
			fmt.Println("服务器关闭连接")
			return
		}
		if err != nil {
			fmt.Println("conn.Read err:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
