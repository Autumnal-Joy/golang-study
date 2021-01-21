package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	go server()

	for i := 0; i < 100; i++ {
		go client(i)
	}

	<-time.After(time.Second * 3)
}

func server() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			return
		}

		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err:", err)
			return
		}

		ret := "---" + string(buf[:n]) + "---"

		n, err = conn.Write([]byte(ret))
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return
		}
		conn.Close()
	}
}

func client(id int) {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Dial err:", err)
	}
	defer conn.Close()

	buf := make([]byte, 4096)

	_, err = conn.Write([]byte(fmt.Sprintf("id: %d", id)))
	if err != nil {
		fmt.Println("conn.Write err:", err)
	}

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Write err:", err)
	}

	fmt.Println(string(buf[:n]))
}
