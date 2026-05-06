package main

import (
	"fmt"
	"log"
	"redirecter/internal/engine"
	"redirecter/internal/packet"
)

func main() {
	fmt.Println("main started")

	e, err := engine.New()
	if err != nil {
		log.Fatalf("error loading Engine: %v", err)
	}

	for {
		pac, err := e.Listen()
		if err != nil {
			fmt.Errorf("main loop listiner error %v", err)
			break
		}

		parsed, err := packet.Parse(&pac)
		if err != nil {
			fmt.Printf("error parsing packet")
		}

		fmt.Printf("Source: %v, Dest: %v", parsed.Source, parsed.Dest)
		if parsed.TCP != nil {
			fmt.Printf(" | TCP Ports: %v -> %v", parsed.TCP.SourcePort, parsed.TCP.DestPort)
		}
		fmt.Println()

		err = e.Send(pac)

		if err != nil {
			log.Printf("błąd wysyłania %v ", err)
		}

	}

	defer e.Close()
}
