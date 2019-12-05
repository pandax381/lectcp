package ip

import (
	"fmt"
	"log"

	"github.com/pandax381/lectcp/pkg/net"
)

type ProtocolRxHandler func(iface net.ProtocolInterface, data []byte, src, dst net.ProtocolAddress) error

type entry struct {
	number    net.ProtocolNumber
	rxHandler ProtocolRxHandler
}

var protocols = map[net.ProtocolNumber]*entry{}

func RegisterProtocol(number net.ProtocolNumber, rxHandler ProtocolRxHandler) error {
	if protocols[number] != nil {
		return fmt.Errorf("protocol `%s` is already registered", number)
	}
	entry := &entry{
		number:    number,
		rxHandler: rxHandler,
	}
	protocols[number] = entry
	log.Printf("protocol registered: %s\n", number)
	return nil
}
