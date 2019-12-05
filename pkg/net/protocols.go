package net

import (
	"fmt"
	"log"
)

type ProtocolRxHandler func(dev *Device, data []byte, src, dst HardwareAddress) error

type packet struct {
	dev  *Device
	data []byte
	src  HardwareAddress
	dst  HardwareAddress
}

type entry struct {
	Type      EthernetType
	rxHandler ProtocolRxHandler
	rxQueue   chan *packet
}

var protocols = map[EthernetType]*entry{}

func RegisterProtocol(Type EthernetType, rxHandler ProtocolRxHandler) error {
	if protocols[Type] != nil {
		return fmt.Errorf("protocol `%s` is registerd", Type)
	}
	entry := &entry{
		Type:      Type,
		rxHandler: rxHandler,
		rxQueue:   make(chan *packet),
	}
	protocols[Type] = entry
	go func() {
		for {
			select {
			case packet, _ := <-entry.rxQueue:
				entry.rxHandler(packet.dev, packet.data, packet.src, packet.dst)
			}
		}
	}()
	log.Printf("protocol registerd: %s\n", entry.Type)
	return nil
}
