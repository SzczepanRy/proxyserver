package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {

	for i := 0; i < 10; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			fmt.Println(err)
		}

		defer conn.Close()

		/*
			"142.250.130.93"
			"443"
		*/
		// 2. Definiujemy DRUGI (ukryty) adres i port, który ma być w BODY
		secretIP := "129.6.15.28"
		secretPort := 13

		// 3. Przygotowujemy nagłówek wewnątrz payloadu (4 bajty na IP + 2 na port = 6 bajtów)
		innerHeader := make([]byte, 6)

		// Konwertujemy string IP na 4 bajty
		parsedIP := net.ParseIP(secretIP).To4()
		copy(innerHeader[0:4], parsedIP)

		// Konwertujemy port int na 2 bajty (Big Endian)
		binary.BigEndian.PutUint16(innerHeader[4:6], uint16(secretPort))

		// 4. (Opcjonalnie) Dodajemy jakąś wiadomość tekstową po adresach
		message := []byte("Hello przez tunel!")

		// Łączymy nagłówek z wiadomością w jeden ostateczny payload
		// [ IP (4B) ] + [ PORT (2B) ] + [ WIADOMOŚĆ ]
		finalPayload := append(innerHeader, message...)

		// 5. Wysyłamy cały pakiet do serwera
		_, err = conn.Write(finalPayload)
		if err != nil {
			fmt.Println("Błąd podczas wysyłania pakietu:", err)
			return
		}

		fmt.Printf("Wysłano pakiet z zaszytym adresem: %s:%d\n", secretIP, secretPort)

		// 4. ODBIERANIE DANYCH W MAIN
		// Ustawiamy timeout (np. 5 sekund), żeby main nie zawiesił się na zawsze,
		// jeśli serwer nic nie odpowie
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		// Tworzymy bufor na odpowiedź od serwera (np. 1024 bajty)
		replyBuffer := make([]byte, 1024)

		// Czytamy odpowiedź z tego samego połączenia
		n, err := conn.Read(replyBuffer)
		if err != nil {
			// Jeśli serwer zamknął połączenie po odesłaniu danych, io.EOF jest normalne
			if err == io.EOF {
				fmt.Println("Serwer zakończył połączenie (EOF).")
			} else {
				fmt.Println("Błąd podczas odbierania odpowiedzi:", err)
				return
			}
		}

		// 5. Wypisujemy to, co serwer nam odpowiedział
		if n > 0 {
			fmt.Printf("Odebrano odpowiedź (%d bajtów):\n", n)

			// 1. Podgląd jako tekst (to co masz teraz)
			fmt.Printf("  Jako tekst: %s\n", string(replyBuffer[:n]))

			// 2. Podgląd surowych wartości bajtów (Liczby od 0 do 255)
			fmt.Printf("  Jako bajty (dec): %v\n", replyBuffer[:n])

			// 3. Podgląd w formacie HEX (idealne do analizy protokołów)
			fmt.Printf("  Jako bajty (hex): %x\n", replyBuffer[:n])
		} else {
			fmt.Println("Połączenie zamknięte, brak danych w odpowiedzi.")
		}

		time.Sleep(1 * time.Second)
	}
}
