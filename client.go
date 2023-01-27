package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("4.查询在线用户")
	fmt.Println("0.退出")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	//_, err := fmt.Scanf("%d", &choice)
	if err != nil {
		fmt.Println(">>>>>>>>请输入合法范围内的数字<<<<<<<<<<<")
		return false
	}

	input = strings.Trim(input, "\r\n")
	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(">>>>>>>>请输入合法范围内的数字<<<<<<<<<<<")
		return false
	}

	if choice >= 0 && choice <= 4 {
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

//查询当前在线用户
func (client *Client) QueryUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err", err)
		return
	}

}

//私聊
func (client *Client) PrivateChat() {
	client.QueryUser()

	var remoteName string
	var chatMsg string
	fmt.Println(">>>>>>>输入你要私聊的对象用户名,exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {

		fmt.Println(">>>>>>>请输出消息,exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "send|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>请输入消息内容, exit退出:")
			fmt.Scanln(&chatMsg)
		}
		client.QueryUser()
		fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
		fmt.Scanln(&remoteName)
	}

}

//公聊
func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>>>请输出消息,exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容,exit退出.")
		fmt.Scanln(&chatMsg)

	}
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
			client.PublicChat()
		case 2:
			//私聊模式
			client.PrivateChat()
		case 3:
			//更新用户名
			client.UpdateName()
		case 4:
			//查询在线用户
			client.QueryUser()
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
