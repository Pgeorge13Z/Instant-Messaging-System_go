//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIP string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIP,
		ServerPort: serverPort,
	}

	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}
	client.conn = conn

	//返回对象
	return client
}

func main() {
	client := NewClient("127.0.0.1", 9999)
	if client == nil {
		fmt.Println(">>>>>>>链接服务器失败....")
		return
	}

	fmt.Println(">>>>>>>>链接服务器成功...")

	//启动客户端业务
	select {}
}
