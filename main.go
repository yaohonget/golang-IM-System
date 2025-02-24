package main

import "fmt"

func main() {
	fmt.Println("Server starting...")
	server := NewServer("", 8888)
	server.Start()
	fmt.Println("Server started...")
}
