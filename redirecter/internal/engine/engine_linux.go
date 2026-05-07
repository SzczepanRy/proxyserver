//go:build linux

package engine

import (
	"fmt"

	"github.com/florianl/go-nfqueue"
)


// linux Potrzebuje  QueueID uint16 pozwala kernelowi rozróżnić, który pakiet wysłać do którego programu.

type LinuxEngine struct {
	nfq *nfqueue.Nfqueue
}

func New() (Engine, error) {
fmt.Println("Inicjalizacja silnika Linux (NFQUEUE)...")

	conf := nfqueue.Config{
		MaxPacketLen: 0xffff, // cały pakietchcemy
		Copymode: nfqueue.NfQnlCopyPacket,


	}

	return &LinuxEngine{}, nil
}

func (e *LinuxEngine) Listen() (Packet, error) {
	// Tu będzie logika NFQUEUE
	return Packet{}, nil
}

func (e *LinuxEngine) Send(p Packet) error {
	return nil
}

func (e *LinuxEngine) Close() {
}
