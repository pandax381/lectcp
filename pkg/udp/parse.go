package udp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"unsafe"

	"github.com/pandax381/lectcp/pkg/net"
)

type header struct {
	SourcePort      uint16
	DestinationPort uint16
	Length          uint16
	Checksum        uint16
}

type datagram struct {
	header
	data []byte
}

func (d datagram) dump() {
	log.Printf("  src port: %d\n", d.SourcePort)
	log.Printf("  dst port: %d\n", d.DestinationPort)
	log.Printf("    length: %d bytes\n", d.Length)
	log.Printf("  checksum: 0x%04x\n", d.Checksum)
	fmt.Println(hex.Dump(d.data))
}

func pseudoHeaderSum(src, dst net.ProtocolAddress, n int) uint32 {
	pseudo := new(bytes.Buffer)
	binary.Write(pseudo, binary.BigEndian, src.Bytes())
	binary.Write(pseudo, binary.BigEndian, dst.Bytes())
	binary.Write(pseudo, binary.BigEndian, uint16(net.ProtocolNumberUDP))
	binary.Write(pseudo, binary.BigEndian, uint16(n))
	return uint32(^net.Cksum16(pseudo.Bytes(), pseudo.Len(), 0))
}

func parse(data []byte, src, dst net.ProtocolAddress) (*datagram, error) {
	hdr := header{}
	if len(data) < int(unsafe.Sizeof(hdr)) {
		return nil, fmt.Errorf("message is too short")
	}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	if int(hdr.Length) > len(data) {
		return nil, fmt.Errorf("length error")
	}
	sum := net.Cksum16(data, len(data), pseudoHeaderSum(src, dst, len(data)))
	if sum != 0 {
		return nil, fmt.Errorf("udp checksum failure: 0x%04x", sum)
	}
	return &datagram{
		header: hdr,
		data:   buf.Bytes(),
	}, nil
}
