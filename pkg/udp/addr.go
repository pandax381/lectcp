package udp

import (
	"fmt"

	"github.com/pandax381/lectcp/pkg/net"
)

type Address struct {
	Addr net.ProtocolAddress
	Port uint16
}

func (a Address) Network() string {
	return "udp"
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Addr, a.Port)
}
