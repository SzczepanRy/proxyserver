//go:build linux

package engine

import (
	"fmt"
)

type LinuxEngine struct {
	// tutaj dodasz pola dla nfqueue
}

func New() (Engine, error) {
	fmt.Println("Inicjalizacja silnika Linux (NFQUEUE)...")
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
