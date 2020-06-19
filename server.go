package main

import (
	"net"
	"os"
	"strconv"
	"time"
)

type Server struct {
	Ip string
	Port uint64
}

func (server *Server) Run() {
	listen, err := net.Listen("tcp", server.Ip+":"+strconv.FormatUint(server.Port, 10))
	if isErrAPrint(err) {os.Exit(1)}
	for {
		conn, err := listen.Accept()
		if isErrAPrint(err) {os.Exit(1)}
		server.handle(&conn)
	}
}

func (server *Server)handle(conn *net.Conn)  {
	defer (*conn).Close()
	if server.checkLogin(conn) {
		server.todoSomething(conn)
	}
}

func (server *Server)checkLogin(conn *net.Conn) bool {
	buffer := make([]byte,512)
	for{
		n, err := (*conn).Read(buffer)
		if isErrAPrint(err) {return false}
		data := make([]byte,n)
		copy(data,buffer[:n+1])
		isRightComplete,err := checkLoginData(data)
		if isErrAPrint(err) {return false}
		if isRightComplete {break}
	}
	return true
}

func (server *Server) todoSomething(conn *net.Conn) {
	go readSomething(conn)
	for{
		writeSomething(conn,[]byte("abc"))
		time.Sleep(time.Second)
	}
}
