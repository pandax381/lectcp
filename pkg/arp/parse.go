package arp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/pandax381/lectcp/pkg/net"
)

type header struct {
	HardwareType          net.HardwareType
	ProtocolType          net.EthernetType
	HardwareAddressLength uint8
	ProtocolAddressLength uint8
	OperationCode         uint16
}

type message struct {
	header
	sourceHardwareAddress []byte
	sourceProtocolAddress []byte
	targetHardwareAddress []byte
	targetProtocolAddress []byte
}

func parse(data []byte) (*message, error) {
	hdr := header{}
	if len(data) < int(unsafe.Sizeof(hdr)) {
		return nil, fmt.Errorf("message is too short")
	}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	msg := message{
		header:                hdr,
		sourceHardwareAddress: make([]byte, hdr.HardwareAddressLength),
		sourceProtocolAddress: make([]byte, hdr.ProtocolAddressLength),
		targetHardwareAddress: make([]byte, hdr.HardwareAddressLength),
		targetProtocolAddress: make([]byte, hdr.ProtocolAddressLength),
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.sourceHardwareAddress); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.sourceProtocolAddress); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.targetHardwareAddress); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.targetProtocolAddress); err != nil {
		return nil, err
	}
	return &msg, nil
}
