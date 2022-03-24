package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := new(User)
	user.Name = userAddr
	user.Addr = userAddr
	user.C = make(chan string)
	user.conn = conn

	go user.ListenMessage()

	return user
}

func (t *User) ListenMessage() {
	msg := <-t.C
	_, err := t.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("conn write err:", err)
	}
}
