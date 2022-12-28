package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Printf("Failed to bind to port 6379: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Printf("Error accepting connection: %v\n", err)
		os.Exit(1)
	}

	done := make(chan struct{})

	go func(c net.Conn) {
		resp := []byte("+PONG\r\n")

		sent, err := c.Write(resp)
		if err != nil {
			fmt.Printf("Error writing to socket: %v\n", err)
		}

		if sent <= 0 {
			fmt.Println("Zero bytes were written")
		}

		c.Close()

		done <- struct{}{}
	}(conn)

	<-done
}
