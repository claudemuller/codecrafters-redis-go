package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
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
	defer c.Close()

	log.Printf("connection from %s accepted", c.RemoteAddr())

	log.Println("attempting to read from connection...")

	buf := make([]byte, 1024)

	for {
		read, err := c.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("EOF")
			} else {
				log.Printf("error reading from socket: %v\n", err)
			}
		}

		log.Printf("%d bytes read\n", read)

		if err := parseResp(c, string(buf)); err != nil {
			log.Println(err.Error())
		}
	}
}

func parseResp(c net.Conn, d string) error {
	tokens := strings.Split(d, "\r\n")

	log.Printf("tokens: %+v\n", tokens)

	for i, v := range tokens {
		dataType := v[:1]
		data := v[1:]

		log.Printf("type: %s\tdata: %s\n", dataType, data)

		switch dataType {
		case "+":
			return handleString(c, data)
		case "$":
			size := data
			data = tokens[i+1]

			log.Printf("type: %s\tdata: %s\tsize: %s\n", dataType, data, size)

			return handleBulkString(c, data, size)
		}
	}

	return nil
}

func handleString(c net.Conn, data string) error {
	if strings.ToLower(data) == "ping" {
		return respondPong(c)
	}

	return nil
}

func handleBulkString(c net.Conn, data, _ string) error {
	if strings.ToLower(data) == "ping" {
		return respondPong(c)
	}

	return nil
}

func respondPong(c net.Conn) error {
	log.Println("respondPong: attempting to write to connection")

	sent, err := c.Write(makeStrResp("PONG"))
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return fmt.Errorf("this is broken pipe error")
		}

		return fmt.Errorf("error writing to socket: %v", err)
	}

	if sent <= 0 {
		return fmt.Errorf("zero bytes were written")
	}

	return nil
}

func makeStrResp(data string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", data))
}
