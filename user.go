package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {
	for msg := range this.C {
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {

	this.server.mapLock.Lock()

	this.server.OnlineMap[this.Name] = this

	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "Online...")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()

	delete(this.server.OnlineMap, this.Name)

	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "Offline...")
}

func (this *User) SendUserMessage(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + " is Online...\n"
			this.SendUserMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendUserMessage("Name is already used..\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendUserMessage("Name changed successfully..\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		targetUserName := strings.Split(msg, "|")[1]
		if targetUserName == "" {
			this.SendUserMessage("Format is incorrect!\n")
			return
		}
		targetUser, ok := this.server.OnlineMap[targetUserName]
		if !ok {
			this.SendUserMessage("The user could not be found..\n")
		} else {
			privateMsg := strings.Split(msg, "|")[2]
			if privateMsg == "" {
				this.SendUserMessage("The messsage is empty.")
			} else {
				targetUser.SendUserMessage(this.Name + " say to you : " + privateMsg + "\n")
			}
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}
