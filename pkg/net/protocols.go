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
		return fmt.Errorf("protocol `%s` is already registered", Type)
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
				if err := entry.rxHandler(packet.dev, packet.data, packet.src, packet.dst); err != nil {
					log.Println(err)
				}
			}
		}
	}()
	log.Printf("protocol registered: %s\n", entry.Type)
	return nil
}
