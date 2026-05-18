package parser

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

func ModifyConn(conn net.Conn) {
	defer conn.Close()
	// timeout na , poaczeniach martwych przez 5 sekund
	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Printf("Błąd ustawiania deadline: %v\n", err)
		return
	}

	header := make([]byte, 6)

	_, err := io.ReadFull(conn, header)
	if err != nil {
		log.Printf("Nie udało się odczytać oryginalnego adresu docelowego: %v\n", err)
		return
	}

	//po otrzymaniu damych wyłączamy timer na połączeniu
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		log.Printf("Błąd resetowania deadline: %v\n", err)
		return
	}

	// wyciąganie ip i poer orginału
	destIP := net.IP(header[:4]).String()
	destPort := binary.BigEndian.Uint16(header[4:])
	originalTarget := net.JoinHostPort(destIP, strconv.Itoa(int(destPort)))

	// połączenie z orginalnym serwerem dest , czeka 5 sekund
	remoteConn, err := net.DialTimeout("tcp", originalTarget, 5*time.Second)
	if err != nil {
		log.Printf("Nie udało się połączyć z oryginalnym celem %s: %v\n", originalTarget, err)
		return
	}else{
		log.Printf("Udało się połączyć z oryginalnym celem %s ", originalTarget)

	}

	defer remoteConn.Close()

	//transparent proxy
	// W tym momencie 'conn' ma już usunięte pierwsze 6 bajtów (skonsumowało je io.ReadFull),
	// więc do oryginalnego serwera trafią już idealnie czyste dane klienta.

	errChan := make(chan error, 2)

	go func() {
		_, err := io.Copy(remoteConn, conn)
		errChan <- err
	}()

	go func() {
		_, err := io.Copy(conn, remoteConn)
		errChan <- err
	}()

	// Czekamy, aż którykolwiek kierunek zakończy transmisję
	<-errChan
}
