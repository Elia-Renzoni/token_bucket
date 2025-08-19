package node

import (
	"net"
	rl "rlimiter/rate_limiter"
	"log"
)

type handler func(net.TCPConn)

type TCPMiddleware struct {
	rlimiter *rl.TokenOwner
}

func InitMiddleware() *TCPMiddleware {
	return &TCPMiddleware{
		rlimiter: rl.InitTokenOwner(),
	}
}

func (t *TCPMiddleware) tokenMiddleware(conn *net.TCPConn, clientProcessor handler) handler {
	return func(conn net.TCPConn) {
		if ok := t.rlimiter.TryTakeToken(); !ok {
			t.drop(conn)
		} else {
			log.Println("Packet Forwarded!")
			clientProcessor(conn)
		}
	}
}

func (t *TCPMiddleware) drop(conn net.TCPConn) {
	conn.Write([]byte("Unvaliable Tokens!"))
	conn.Close()
	log.Println("Dropped Packet!")
}
