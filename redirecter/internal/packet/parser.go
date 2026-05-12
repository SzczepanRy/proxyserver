package packet

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"redirecter/internal/engine"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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

// tutaj dodam ipBytes i portBytes
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

func Modify(p *engine.Packet, ipBytes []byte, portBytes []byte) (*engine.Packet, error) {
	// dekodujemy pakiet na warstwy
	packet := gopacket.NewPacket(p.Data, layers.LayerTypeIPv4, gopacket.Default)

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)

	if ipLayer == nil || tcpLayer == nil {
		return p, nil //  puszczamy bez zmian bo nie tcp
	}

	ip, _ := ipLayer.(*layers.IPv4)
	tcp, _ := tcpLayer.(*layers.TCP)

	oldDestIP := make([]byte, 4)
	copy(oldDestIP, ip.DstIP)

	oldDestPort := make([]byte, 2)
	binary.BigEndian.PutUint16(oldDestPort, uint16(tcp.DstPort))

	marker := append(oldDestIP, oldDestPort...)

	ip.DstIP = net.IP(ipBytes)
	tcp.DstPort = layers.TCPPort(binary.BigEndian.Uint16(portBytes))

	// Dodajemy 6 bajtów na początku istniejących danych
	tcp.Payload = append(marker, tcp.Payload...)

	// gopacket liczy sumy kontrolne
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true, // Oblicz sumy kontrolne (IP i TCP!)
		FixLengths:       true, // Popraw Total Length w IP
	}

	// TCP Checksum wymaga pseudo-nagłówka IP,
	// dlatego wywołujemy to specjalne polecenie:
	err := tcp.SetNetworkLayerForChecksum(ip)
	if err != nil {
		return nil, fmt.Errorf("błąd checksumy: %v", err)
	}

	// Składamy wszystko z powrotem do []byte
	err = gopacket.SerializeLayers(buffer, options, ip, tcp, gopacket.Payload(tcp.Payload))
	if err != nil {
		return nil, fmt.Errorf("błąd serializacji: %v", err)
	}

	p.Data = buffer.Bytes()
	return p, nil
}

func calculateIPChecksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data); i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}
	for sum > 0xffff {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return uint16(^sum)
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
