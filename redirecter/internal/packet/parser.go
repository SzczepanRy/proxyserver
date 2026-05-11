package packet

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"redirecter/internal/engine"
	"strconv"
)

type TCPHeader struct {
	SourcePort uint16
	DestPort   uint16
}

type ParsedPacket struct {
	Raw      *engine.Packet
	Source   net.IP
	Dest     net.IP
	Protocol uint8
	TCP      *TCPHeader
}

//tutaj dodam ipBytes i portBytes
// przez co było by najlepiej
func Parse(p *engine.Packet) (*ParsedPacket, error) {

	data := p.Data
	if len(data) < 20 {
		return nil, fmt.Errorf("pakiet za kturki")
	}

	srcIP := net.IP(data[12:16])
	destIP := net.IP(data[16:20])
	proto := data[9]

	parsed := &ParsedPacket{
		Raw:      p,
		Source:   srcIP,
		Dest:     destIP,
		Protocol: proto,
	}

	// proto  6 == Tcp

	if proto == 6 && len(data) >= 40 {
		paloadOffset := (data[0] & 0x0F) * 4 // 4 len IP
		tcpData := data[paloadOffset:]
		parsed.TCP = &TCPHeader{
			SourcePort: binary.BigEndian.Uint16(tcpData[0:2]),
			DestPort:   binary.BigEndian.Uint16(tcpData[2:4]),
		}

	}

	return parsed, nil
}

func LoadEnv() (ipBytes []byte, portBytes []byte, err error) {
	ipStr := os.Getenv("PROXYIP")
	if ipStr == "" {
		return nil, nil, fmt.Errorf("PARSER ERROR : nie znaleziono pola PROXYIP")
	}
	ip := net.ParseIP(ipStr)
	//conv to 4 bytes
	ipBytes = ip.To4()
	if ipBytes == nil {
		return nil, nil, fmt.Errorf("PARSER ERROR : bład zamiany IP na bajty")
	}

	portStr := os.Getenv("PROXYPORT")
	if portStr == "" {
		return nil, nil, fmt.Errorf("PARSER ERROR : nie znaleziono pola PROXYPORT")
	}

	port, _ := strconv.Atoi(portStr)
	portUint16 := uint16(port)
	portBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, portUint16)

	return ipBytes, portBytes, nil

}
