package main

import (
	"rlimiter/node"
	"log"
)

func startRoutines(s node.Service) {
	go s.Bind()
	go s.Serve()
}

func main() {
	waiter := make(chan struct{})
	server := node.InitNode("127.0.0.1", "6060", waiter)
	startRoutines(server)
	<- waiter
	log.Printf("Server Shutdown!")
}
