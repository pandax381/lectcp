package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/pandax381/lectcp/pkg/net"
)

const IPVersion4 = 4

type header struct {
	VHL      uint8
	TOS      uint8
	Len      uint16
	Id       uint16
	Offset   uint16
	TTL      uint8
	Protocol net.ProtocolNumber
	Sum      uint16
	Src      Address
	Dst      Address
}

type datagram struct {
	header
	payload []byte
}

func parse(data []byte) (*datagram, error) {
	hdr := header{}
	if len(data) < int(unsafe.Sizeof(hdr)) {
		return nil, fmt.Errorf("ip packet is too short (%d)", len(data))
	}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	if hdr.VHL>>4 != IPVersion4 {
		return nil, fmt.Errorf("not ipv4 packet")
	}
	hlen := int((hdr.VHL & 0x0f) << 2)
	if len(data) < hlen {
		return nil, fmt.Errorf("need least header length's data")
	}
	sum := net.Cksum16(data, hlen, 0)
	if sum != 0 {
		return nil, fmt.Errorf("ip checksum error (%x)", sum)
	}
	if len(data) < int(hdr.Len) {
		return nil, fmt.Errorf("ip packet length error")
	}
	if hdr.TTL == 0 {
		return nil, fmt.Errorf("ip packet was dead (TTL=0)")
	}
	return &datagram{
		header:  hdr,
		payload: data[hlen:int(hdr.Len)],
	}, nil
}
