package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User //key:string,value:*User
	mapLock   sync.RWMutex
	//全局消息message
	Message chan string
}

func (t *Server) ListenMessage() {
	for {
		msg := <-t.Message
		t.mapLock.Lock()
		for _, cli := range t.OnlineMap {
			cli.C <- msg
		}
		t.mapLock.Unlock()
	}
}

func (t *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Name + "]" + user.Addr + msg
	t.Message <- sendMsg
}

func (t *Server) Handler(conn net.Conn) {
	//fmt.Println("链接建立成功！")
	user := NewUser(conn)
	t.mapLock.Lock()
	t.OnlineMap[user.Name] = user //why username?
	t.mapLock.Unlock()
	t.BroadCast(user, "已上线")
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				t.BroadCast(user, "已下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			msg := string(buf[:n-1])
			t.BroadCast(user, msg)
		}
	}()

	select {}
}

func NewServer(ip string, port int) *Server {
	server := new(Server)
	server.Ip = ip
	server.Port = port
	server.OnlineMap = make(map[string]*User)
	server.Message = make(chan string)
	return server
}

func (t *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.Ip, t.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}
	defer listener.Close()
	go t.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		go t.Handler(conn)
	}
}
