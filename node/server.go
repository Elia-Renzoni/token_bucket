package node

import (
	"errors"
	"io"
	"log"
	"net"
)

type Service interface {
	Bind()
	Serve()
}

type Node struct {
	ls            *net.TCPListener
	addr          *net.TCPAddr
	waiter        chan struct{}
	acceptor      chan struct{}
	host, port    string
	networkError  error
	tcpMiddleware *TCPMiddleware
}

func InitNode(host, port string, waiter chan struct{}) *Node {
	complete, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return &Node{
		addr:          complete,
		waiter:        waiter,
		acceptor:      make(chan struct{}),
		host:          host,
		port:          port,
		tcpMiddleware: InitMiddleware(),
	}
}

func (n *Node) Bind() {
	n.ls, n.networkError = net.ListenTCP("tcp", n.addr)
	if n.networkError != nil {
		n.waiter <- struct{}{}
		log.Fatal(n.networkError.Error())
		return
	}
	n.acceptor <- struct{}{}
}

func (n *Node) Serve() {
	<-n.acceptor
	log.Println("Server Listening...")

	for {
		tcpConn, err := n.ls.AcceptTCP()
		if err != nil {
			n.waiter <- struct{}{}
			log.Fatal(n.networkError.Error())
			break
		}
		sessionHandler := n.tcpMiddleware.tokenMiddleware(tcpConn, n.handleClient)
		sessionHandler(*tcpConn)
	}
}

func (n *Node) handleClient(client net.TCPConn) {
	var (
		length int
		err    error
	)
	defer client.Close()
	buf := make([]byte, 2024)

	for {
		length, err = client.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("read all the data from the socket")
				break
			}
			break
		}
	}

	data := string(buf[:length])
	if data == "ping" {
		client.Write([]byte("pong"))
	}
}
