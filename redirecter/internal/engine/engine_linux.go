//go:build linux

package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/florianl/go-nfqueue"
	"github.com/mdlayher/netlink"
)

///"golang.org/x/sys/unix"

type LinuxEngine struct {
	nfq     *nfqueue.Nfqueue
	packets chan Packet        // kanał miedzy callback a listen
	cancel  context.CancelFunc // zapisujemy cancel, aby użyć w Close()
}

func New() (Engine, error) {
	fmt.Println("Inicjalizacja silnika Linux (NFQUEUE)...")

	conf := nfqueue.Config{
		NfQueue:      100,    // numer kolejki iptables
		MaxPacketLen: 0xFFFF, // cały pakietchcemy
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		AfFamily:     0, // rodzina ip4
		//AfFamily:     unix.AF_INET, // rodzina ip4
		//ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
	}

	nf, err := nfqueue.Open(&conf)

	if err != nil {
		return nil, fmt.Errorf("Error loading nfqueue engine %v", err)
	}

	// Avoid receiving ENOBUFS errors.
	if err := nf.SetOption(netlink.NoENOBUFS, true); err != nil {
		return nil, fmt.Errorf("failed to set netlink option %v: %v\n", netlink.NoENOBUFS, err)
	}

	// Tworzymy kontekst główny dla biblioteki nfqueue
	ctx, cancel := context.WithCancel(context.Background())

	e := &LinuxEngine{nfq: nf, packets: make(chan Packet, 1024), cancel: cancel}

	handlePacket := func(p nfqueue.Attribute) int {

		fmt.Printf("[NFQUEUE] Przechwycono pakiet ID: %d, Rozmiar: %d\n", *p.PacketID, len(*p.Payload))
		if p.Payload != nil && p.PacketID != nil {
			// this sould be better

			e.packets <- Packet{
				Data:    *p.Payload,
				Context: *p.PacketID, // tu zapisujemy uint32 wymagany przez SetVerdict
			}
			//tymczasowe
			//nf.SetVerdict(*p.PacketID, nfqueue.NfAccept)
		}
		return 0
	}

	handleError := func(err error) int {
		fmt.Printf("NFQUEUE błąd tła: %v\n", err)
		return 0
	}

	// Odpalamy rejestrację - pamiętaj, że to NIE blokuje!
	err = nf.RegisterWithErrorFunc(ctx, handlePacket, handleError)
	if err != nil {
		// Jeśli sama rejestracja się nie powiodła, sprzątamy i rzucamy błąd
		cancel()
		nf.Close()
		return nil, fmt.Errorf("FATAL: RegisterWithErrorFunc zwróciło błąd: %v", err)
	}

	fmt.Println("Silnik NFQUEUE działa i nasłuchuje.")
	return e, nil
}

func (e *LinuxEngine) Listen() (Packet, error) {
	p, ok := <-e.packets
	if !ok {
		return Packet{}, fmt.Errorf("engine close")
	}

	return p, nil
}

func (e *LinuxEngine) Send(p Packet) error {

	packetID, ok := p.Context.(uint32)
	if !ok {
		return fmt.Errorf("invalidpacket context could not find packetID ")
	}

	//wydanie werdyktu ACCEPT1 drop 0
	// p.Data zawiera Twoje zmodyfikowane bajty
	return e.nfq.SetVerdictModPacket(packetID, 1, p.Data)

}

func (e *LinuxEngine) Close() {
	// Najpierw anulujemy kontekst, co zatrzymuje procesy biblioteki
	if e.cancel != nil {
		e.cancel()
	}
	// Zamykamy gniazdo
	if e.nfq != nil {
		e.nfq.Close()
	}
	// Zamykamy kanał, żeby Listen() mogło bezpiecznie zwrócić błąd i się zakończyć
	close(e.packets)
}
