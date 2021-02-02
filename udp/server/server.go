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
	// 解析地址
	srvAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.ResolveUDPAddr err: ", err)
		return
	}
	// 启动监听
	udpConn, err := net.ListenUDP("udp", srvAddr)
	if err != nil {
		fmt.Println("net.ListenUDP err: ", err)
		return
	}
	defer udpConn.Close()
	fmt.Println("服务器已就绪")

	buf := make([]byte, 4096)
	for {
		// 建立连接
		n, cltAddr, err := udpConn.ReadFromUDP(buf)
		go func() {
			if n == 0 {
				fmt.Println("客户端关闭连接")
				return
			}
			if err != nil {
				fmt.Println("udpConn.ReadFromUDP err:", err)
				return
			}

			// 数据处理
			str := strings.TrimSpace(string(buf[:n]))
			fmt.Printf("客户端: %v 发送数据: %q\n", cltAddr, str)
			res := "respond: " + strings.ToUpper(str)

			// 数据回写
			n, err = udpConn.WriteToUDP([]byte(res), cltAddr)
			if err != nil {
				fmt.Println("udpConn.WriteToUDP err:", err)
				return
			}
		}()
	}
}
