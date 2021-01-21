package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// 用户结构
type Client struct {
	// 用户消息通道
	c    chan string
	name string
	// 用户地址ip+端口
	addr string
}

// 全局消息通道
var message = make(chan string)

// 在线用户表
var onlineUserMap = make(map[string]Client)

func main() {
	server()
}

func server() {
	// 服务器绑定IP+端口
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	// 创建广播消息管理者go程
	go manager()

	fmt.Println("服务器就绪")

	for {
		// 等待用户连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			return
		}
		// 分配go程处理用户连接
		go handleConnect(conn)
	}
}

// 广播消息管理go程
func manager() {
	// 监听用户消息
	for {
		msg := <-message

		// 广播消息给所有用户
		for _, client := range onlineUserMap {
			client.c <- msg
		}
	}
}

// 用户连接处理go程
func handleConnect(conn net.Conn) {
	defer conn.Close()

	// 获取用户地址网络地址
	addr := conn.RemoteAddr().String()
	fmt.Println("用户连接", addr)

	// 创建新用户, 用户名默认为地址
	client := Client{make(chan string), addr, addr}

	// 将新用户添加到用户表中
	onlineUserMap[addr] = client

	// 创建响应用户的go程
	go writeMsgToClient(conn, client)

	// 发送用户上线通知到全局消息通道message中
	message <- wrapMsg(client, "login")

	// 判断用户退出
	quit := make(chan bool)
	// 判断用户活跃 make(chan bool)
	isActive := make(chan bool)

	// 处理用户发送的消息
	go handleMessage(conn, client, quit, isActive)

	for {
		select {
		case <-quit:
			// 用户主动退出
			delete(onlineUserMap, client.addr)
			message <- wrapMsg(client, "logout")
			return
		case <-time.After(time.Second * 10):
			// 超时强制踢出
			delete(onlineUserMap, client.addr)
			message <- wrapMsg(client, "logout")
			return
		case <-isActive:
			// 重置计时器
		}
	}
}

// 用户消息响应go程
func writeMsgToClient(conn net.Conn, client Client) {
	// 监听用户消息通道是否有消息
	for msg := range client.c {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("writeMsgToClient err: ", err)
		}
	}
}

// 包装响应给用户的信息
func wrapMsg(client Client, msg string) (buf string) {
	buf = "[" + client.addr + "] " + client.name + ": " + msg
	return
}

// 解析用户发送信息, 并相应处理
func handleMessage(conn net.Conn, client Client, quit chan<- bool, isActive chan<- bool) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			quit <- true
			return
		}
		if err != nil {
			fmt.Println("handleMessage err: ", err)
			quit <- true
			return
		}
		msg := strings.TrimSpace(string(buf[:n]))
		fmt.Printf("用户%s输入: %q\n", client.name, msg)

		if msg == "/who" {
			// 用户查询在线列表
			userList := "online user list:\n"
			for _, user := range onlineUserMap {
				userList += user.addr + ": " + user.name + "\n"
			}
			client.c <- userList
		} else if strings.HasPrefix(msg, "/rename ") {
			// 用户改名
			client.name = msg[8:]
			onlineUserMap[client.addr] = client
		} else if msg == "/quit" {
			quit <- true
		} else {
			// 发送用户消息到全局消息通道
			message <- wrapMsg(client, msg)
		}
		isActive <- true
	}
}
