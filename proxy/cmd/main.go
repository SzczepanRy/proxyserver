package main

import (
	"fmt"
	"log"
	"net"
	"proxy/internal/parser"
)

func main() {

	l, err := net.Listen("tcp", ":8080")

	fmt.Println("proxy started")

	if err != nil {
		log.Fatalf("could not init Listener , %v", err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Printf("could not init Listener , %v", err)
			continue
		}

		go parser.ModifyConn(conn)

	}

}
