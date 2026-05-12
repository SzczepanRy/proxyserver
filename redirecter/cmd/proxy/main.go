package main

import (
	"fmt"
	"log"
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

	ipBytes, portBytes, err := packet.LoadEnv()

	if err != nil {
		log.Fatalf("Błąd ładowania zmiennych pliku .env : %v", err)
	}

	for {
		pac, err := e.Listen()

		if err != nil {
			log.Printf("main loop listiner error %v", err)
			break
		}

		// czysto kosmetyczne

		parsed, err := packet.Parse(&pac)
		if err != nil {
			fmt.Printf("error parsing packet : %v" , err )
		}

		if parsed.TCP != nil {
			fmt.Printf("Source: %v, Dest: %v", parsed.Source, parsed.Dest)
			fmt.Printf(" | TCP Ports: %v -> %v  \n ", parsed.TCP.SourcePort, parsed.TCP.DestPort)
		}
		// end
		np, err := packet.Modify(&pac, ipBytes, portBytes)

		if err != nil {
			log.Fatalf("could not Modify packet : %v ", err)
		}


		parsed, err = packet.Parse(np)
		if err != nil {
			fmt.Printf("error parsing packet : %v" , err )
		}

		if parsed.TCP != nil {
			fmt.Printf("Source: %v, Dest: %v", parsed.Source, parsed.Dest)
			fmt.Printf(" | TCP Ports: %v -> %v  \n ", parsed.TCP.SourcePort, parsed.TCP.DestPort)
		}



		if np != nil {
			err = e.Send(*np)
			if err != nil {
				log.Printf("błąd wysyłania %v ", err)
			}
		}

	}

	defer e.Close()
}
