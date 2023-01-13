//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	choice     int
}

func NewClient(serverIP string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIP,
		ServerPort: serverPort,
		choice:     333,
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

func (Client *Client) menu() bool {
	var choice int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&choice)

	if choice >= 0 && choice <= 3 {
		Client.choice = choice
		return true
	} else {
		fmt.Println(">>>>>>>>请输入合法范围内的数字<<<<<<<<<")
		return false
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
	//io.copy等于下面的功能
	// for {
	// 	buf := make([]byte, 4096)
	// 	client.conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.choice != 0 {
		for client.menu() != true {
		}
		switch client.choice {
		case 1:
			//公聊模式
			fmt.Println("公聊")
		case 2:
			//私聊模式
			fmt.Println("私聊")
		case 3:
			//更新用户名
			client.UpdateName()
		}

	}

}

//定义命令行解析的变量
var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器的IP地址,默认是127.0.0.1")
	flag.IntVar(&serverPort, "port", 9999, "设置服务器的端口,默认是9999")
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>链接服务器失败....")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>>>>>>链接服务器成功...")

	//启动客户端业务
	client.Run()
}
