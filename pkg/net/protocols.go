package net

import (
	"fmt"
	"log"
	"sync"
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

var protocols = sync.Map{}

func RegisterProtocol(Type EthernetType, rxHandler ProtocolRxHandler) error {
	if _, exists := protocols.Load(Type); exists {
		return fmt.Errorf("protocol `%s` is already registered", Type)
	}
	entry := &entry{
		Type:      Type,
		rxHandler: rxHandler,
		rxQueue:   make(chan *packet),
	}
	// launch rx loop
	go func() {
		for packet := range entry.rxQueue {
			if err := entry.rxHandler(packet.dev, packet.data, packet.src, packet.dst); err != nil {
				log.Println(err)
			}
		}
	}()
	protocols.Store(Type, entry)
	log.Printf("protocol registered: %s\n", entry.Type)
	return nil
}
