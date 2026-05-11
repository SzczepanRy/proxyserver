package main

import (
	"fmt"
	"log"
	"os"
	"redirecter/internal/engine"
	"redirecter/internal/packet"

	"github.com/joho/godotenv"
)

func main() {

	e, err := engine.New()
	if err != nil {
		log.Fatalf("error loading Engine: %v", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Błąd ładowania pliku .env")
	}

	ipBytes , portBytes , err := packet.LoadEnv()

	if err != nil {
		log.Fatalf("Błąd ładowania zmiennych pliku .env : %v" , err)
	}

	for {
		pac, err := e.Listen()

		if err != nil {
			log.Printf("main loop listiner error %v", err)
			break
		}

		fmt.Fprintf(os.Stderr, ">>> Odebrano pakiet! ID: %v | Rozmiar: %d\n", pac.Context, len(pac.Data))

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
