package main

import (
	"net"
	"log"
	"sync"
)

func main() {
	address := net.JoinHostPort("127.0.0.1", "6060")
	var group sync.WaitGroup
	group.Add(300)
	
	for i := 0; i < 300; i++ {
		go func() {
			defer group.Done()
			conn, err := net.Dial("tcp", address)
			if err != nil {
				log.Println("Error in Dialer Function!")
				return
			}

			defer conn.Close()

			conn.Write([]byte("ping"))
			buffer := make([]byte, 1024)

			var (
				ln int
				errn error
			)
			ln, errn = conn.Read(buffer)
			if errn != nil {
				return
			}

			log.Printf("Ack: %s \n", string(buffer[:ln]))
		}()
	}
	group.Wait()
}
