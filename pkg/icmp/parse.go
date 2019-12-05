package icmp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"unsafe"

	"github.com/pandax381/lectcp/pkg/net"
)

type messageType uint8

const (
	messageTypeEchoReply messageType = 0
	messageTypeEcho      messageType = 8
)

func (m messageType) String() string {
	switch m {
	case messageTypeEchoReply:
		return "Echo Reply"
	case messageTypeEcho:
		return "Echo"
	default:
		return "Unknown"
	}
}

type header struct {
	Type messageType
	Code uint8
	Sum  uint16
}

type message interface {
	messageType() messageType
	marshal() ([]byte, error)
	dump()
}

func parse(data []byte) (message, error) {
	hdr := header{}
	if len(data) < int(unsafe.Sizeof(hdr)) {
		return nil, fmt.Errorf("message is too short")
	}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	switch hdr.Type {
	case messageTypeEcho, messageTypeEchoReply:
		message := messageEcho{header: hdr}
		if err := binary.Read(buf, binary.BigEndian, &message.id); err != nil {
			return nil, err
		}
		if err := binary.Read(buf, binary.BigEndian, &message.seq); err != nil {
			return nil, err
		}
		message.data = buf.Bytes()
		return &message, nil
	default:
		message := messageGeneral{header: hdr}
		message.data = buf.Bytes()
		return &message, nil
	}
}

type messageGeneral struct {
	header
	data []byte
}

func (m *messageGeneral) messageType() messageType {
	return m.Type
}

func (m *messageGeneral) marshal() ([]byte, error) {
	buf := make([]byte, 4+len(m.data))
	buf[0] = uint8(m.Type)
	buf[1] = uint8(m.Code)
	binary.BigEndian.PutUint16(buf[2:4], 0)
	copy(buf[4:], m.data)
	sum := net.Cksum16(buf, len(buf), 0)
	binary.BigEndian.PutUint16(buf[2:4], sum)
	return buf, nil
}

func (m *messageGeneral) dump() {
	fmt.Printf("  type: %s (%x)\n", m.Type, uint8(m.Type))
	fmt.Printf("  code: 0x%02x\n", uint8(m.Code))
	fmt.Printf("   sum: 0x%04x\n", m.Sum)
	fmt.Printf("  data: %d bytes\n", len(m.data))
	fmt.Printf("%s", hex.Dump(m.data))
}

type messageEcho struct {
	header
	id   uint16
	seq  uint16
	data []byte
}

func (m *messageEcho) messageType() messageType {
	return m.Type
}

func (m *messageEcho) marshal() ([]byte, error) {
	buf := make([]byte, 4+4+len(m.data))
	buf[0] = uint8(m.Type)
	buf[1] = uint8(m.Code)
	binary.BigEndian.PutUint16(buf[2:4], 0)
	binary.BigEndian.PutUint16(buf[4:6], uint16(m.id))
	binary.BigEndian.PutUint16(buf[6:8], uint16(m.seq))
	copy(buf[8:], m.data)
	sum := net.Cksum16(buf, len(buf), 0)
	binary.BigEndian.PutUint16(buf[2:4], sum)
	return buf, nil
}

func (m *messageEcho) dump() {
	fmt.Printf("  type: %s (%x)\n", m.Type, uint8(m.Type))
	fmt.Printf("  code: 0x%02x\n", uint8(m.Code))
	fmt.Printf("   sum: 0x%04x\n", m.Sum)
	fmt.Printf("    id: 0x%04x\n", m.id)
	fmt.Printf("   seq: 0x%04x\n", m.seq)
	fmt.Printf("  data: %d bytes\n", len(m.data))
	fmt.Printf("%s", hex.Dump(m.data))
}
