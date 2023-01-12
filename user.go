package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn

	server *Server
}

//创建获取用户的API
func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		C:      make(chan string),
		Conn:   conn,
		server: server,
	}

	//驱动监听当前user channel消息的goroutine
	go user.Listen_userMessage()

	return user

}

func (this *User) Online() {
	//用户上线,将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {

	//用户下线,将用户从onlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "下线")

}

func (this *User) SendMsg(msg string) {
	this.Conn.Write([]byte(msg))
}

//消息处理
func (this *User) Domsg(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[0:7] == "rename|" {
		newname := strings.Split(msg, "|")[1]

		//判断name是否存在
		_, ok := this.server.OnlineMap[newname]
		if ok {
			this.SendMsg("当前用户名已经被使用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newname] = this
			this.server.mapLock.Unlock()

			this.Name = newname
			this.SendMsg("您已经更新用户名:" + newname + "\n")
		}
	} else if len(msg) > 6 && msg[0:5] == "send|" {
		remotename := strings.Split(msg, "|")[1]
		if remotename == "" {
			this.SendMsg("消息格式不正确,请使用\"send|张三|你好\"的格式")
			return
		} else if remotename == this.Name {
			this.SendMsg("您不能对自己发送消息")
		} else {
			remoteUser, ok := this.server.OnlineMap[remotename]
			if !ok {
				this.SendMsg("该用户不存在\n")
				return
			}

			content := strings.Split(msg, "|")[2]
			if content == "" {
				this.SendMsg("无消息内容，请重发\n")
				return
			}
			remoteUser.SendMsg(this.Name + "对您说:" + content)
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

//监听当前User channel方法，一旦有消息，直接发送给客户端
func (this *User) Listen_userMessage() {
	for {
		msg := <-this.C
		this.SendMsg(msg + "\n")
	}
}
