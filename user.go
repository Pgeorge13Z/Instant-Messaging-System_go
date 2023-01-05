package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
}

//创建获取用户的API
func NewUser(conn net.Conn) *User {
	user := &User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		C:    make(chan string),
		Conn: conn,
	}

	//驱动监听当前user channel消息的goroutine
	go user.Listen_userMessage()

	return user

}

//监听当前User channel方法，一旦有消息，直接发送给客户端
func (this *User) Listen_userMessage() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg + "\n"))
	}
}
