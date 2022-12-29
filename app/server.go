package main

import (
	"log"
	"net"
	"os"
)

func main() {
	log.SetPrefix("rediz: ")

	const port = "6379"

	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Printf("failed to bind to port 6379: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	log.Printf("server started on port %s", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v\n", err)
			os.Exit(1)
		}

		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	log.Printf("connection from %s accepted", c.RemoteAddr())

	resp := []byte("+PONG\r\n")

	sent, err := c.Write(resp)
	if err != nil {
		log.Printf("error writing to socket: %v\n", err)
	}

	if sent <= 0 {
		log.Println("zero bytes were written")
	}

	c.Close()
}
