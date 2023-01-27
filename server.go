package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

//监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) Listen_serverMessager() {
	for {
		msg := <-this.Message

		//将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

//广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	//fmt.Println("链接建立成功")

	user := NewUser(conn, this)

	//用户上线
	user.Online()

	isalive := make(chan bool)

	//接受客户端发送的消息
	go func() {

		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {

				//用户下线
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println(("Conn read erro:"), err)
				return
			}

			//提取用户的消息（去除'\n')
			msg := string(buf[:n-1])
			//msg := string(buf)

			//将得到的进行处理
			user.Domsg(msg)

			//更新alive
			isalive <- true
		}
	}()

	//当前handler阻塞
	//超时强踢功能
	for {
		select {
		case <-isalive:
		//当前用户活跃，重置定时器
		//不做事，进入下一次循环

		case <-time.After(time.Second * 300):
			//超时，关闭用户
			user.SendMsg("你被踢了")
			//关闭user使用的管道
			close(user.C)
			//关闭conn资源
			conn.Close()

			//退出当前handler
			return
		}
	}

}

//启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go this.Listen_serverMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}
