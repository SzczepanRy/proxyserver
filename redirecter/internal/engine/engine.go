package engine

type Packet struct {
	Data    []byte
	Context interface{}
}

type Engine interface {
	Listen() (Packet, error)
	Send(Packet) error
	Close()
}
