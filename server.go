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

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	fmt.Println("Server created...")
	return server
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection successful.")

	user := NewUser(conn)

	this.mapLock.Lock()

	this.OnlineMap[user.Name] = user

	this.mapLock.Unlock()

	this.BroadCast(user, "Online...")

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.BroadCast(user, "Offline...")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:n-1])
			this.BroadCast(user, msg)
		}
	}()

	select {}

}

func (this *Server) Start() {
	fmt.Println("Connection establishing, ip : ", this.Ip, ", port : ", this.Port, ", connection string : ", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	//listener, err := net.Listen("tcp", "127.0.0.1:8082")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err : ", err)
		return
	}
	defer listener.Close()

	go this.ListenMessage()

	fmt.Println("Actually listening on", listener.Addr())
	if listener == nil {
		fmt.Println("listener is empty.")
		return
	} else {
		fmt.Printf("listener is %v\n", listener)
	}
	for {
		fmt.Println("starting listener accept()")
		conn, err := listener.Accept()
		fmt.Println("listener accepted...")
		if err != nil {
			fmt.Println("listener accept err :", err)
			continue
		}

		go this.Handler(conn)
	}
	/*
		tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", this.Ip, this.Port))
		if err != nil {
			fmt.Println("net.ResolveTCPAddr err : ", err)
			os.Exit(-1)
		}
		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			fmt.Println("net.ListenTcp err : ", err)
			os.Exit(-1)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			fmt.Println("Accept a new connection")
			go this.Handler(conn)
		}*/
	/*l, err := net.Listen("tcp", ":2000")
	if err != nil {
		fmt.Println("listener listen err :", err)
	}
	defer l.Close()
	fmt.Println("Actually listening on", l.Addr())
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("listener accept err :", err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}*/
}
