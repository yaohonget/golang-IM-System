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
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err : ", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {

	var flag int

	fmt.Println("1. Boardcast")
	fmt.Println("2. Privat Message")
	fmt.Println("3. Rename Username")
	fmt.Println("0. Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>> Please input the correct number. <<<<<<<")
		return false
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println("Please input your new name.")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Conn.Write err : ", err)
		return false
	}
	return true
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("Please input your message, exit by input 'exit'.")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("Conn.Write err : ", err)
				break
			}

		}

		chatMsg = ""
		fmt.Println("Please input your message, exit by input 'exit'.")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Conn.Write err : ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	client.SelectUsers()
	fmt.Println("Please input the target username, exit by input 'exit'.")
	fmt.Scanln(remoteName)
	for remoteName != "exit" {
		fmt.Println("Please input your message, exit by input 'exit'.")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("Conn.Write err : ", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("Please input your message, exit by input 'exit'.")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println("Please input the target username, exit by input 'exit'.")
		fmt.Scanln(remoteName)
	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			fmt.Println("Boardcast selected...")
			client.PublicChat()
			break
		case 2:
			fmt.Println("Private Message selected...")
			client.PrivateChat()
			break
		case 3:
			fmt.Println("Rename selected...")
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Setup IP, default is 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "Setup IP, default is 8888")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>> Connection failed <<<<<<<<<")
		return
	}

	go client.DealResponse()
	fmt.Println(">>>>>>> Connection successful <<<<<<<<<")

	client.Run()
}
