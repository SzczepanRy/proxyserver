//go:build windows

package engine

import (
	"fmt"

	"github.com/imgk/divert-go"
)

type WindowsEngine struct {
	handle *divert.Handle
}

func New() (Engine, error) {
	h, err := divert.Open("ip", divert.LayerNetwork, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("WindDivert ERROR: %v" ,err)
	}

	return &WindowsEngine{
		handle: h,
	}, nil

}

func (e *WindowsEngine) Listen() (Packet, error) {
	buf := make([]byte, 1500)//surowe dane pakiety , std MTU 1500
	addr := &divert.Address{}

	n, err := e.handle.Recv(buf, addr)

	if err != nil {
		return Packet{}, fmt.Errorf("błąd odbierania obiktu : %w ", err)
	}

	// addr zwiera metadane , do odesłania adresu
	return Packet{
		Data:    buf[:n],
		Context: addr,
	}, nil

}

func (e *WindowsEngine) Send(p Packet) error {
	// tu wyciąga addr z kontexut
	addr , ok :=  p.Context.(*divert.Address)

	if !ok {
		return fmt.Errorf("nie dało rady wyciądnąc addr z ctx ")
	}

	//przesyłam dalej
	_, err := e.handle.Send(p.Data, addr)

	return err

}

func (e *WindowsEngine) Close() {
	if e.handle != nil {
		e.handle.Close()
	}
}
