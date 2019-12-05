package ip

import (
	"github.com/pandax381/lectcp/pkg/net"
)

type Interface struct {
	unicast   Address
	netmask   Address
	broadcast Address
	gateway   Address
	device    *net.Device
}

func newInterface(dev *net.Device, unicast, netmask Address) (*Interface, error) {
	return &Interface{
		unicast: unicast,
		netmask: netmask,
		broadcast: Address{
			unicast[0]&netmask[0] | ^netmask[0],
			unicast[1]&netmask[1] | ^netmask[1],
			unicast[2]&netmask[2] | ^netmask[2],
			unicast[3]&netmask[3] | ^netmask[3],
		},
		device: dev,
	}, nil
}

func CreateInterface(dev *net.Device, unicast, netmask, gateway string) (*Interface, error) {
	addr, err := ParseAddress(unicast)
	if err != nil {
		return nil, err
	}
	mask, err := ParseAddress(netmask)
	if err != nil {
		return nil, err
	}
	iface, err := newInterface(dev, addr, mask)
	if err != nil {
		return nil, err
	}
	return iface, nil
}

func (iface *Interface) Type() net.EthernetType {
	return net.EthernetTypeIP
}

func (iface *Interface) Address() net.ProtocolAddress {
	return iface.unicast
}

func (iface *Interface) Device() *net.Device {
	return iface.device
}

func (iface *Interface) Tx(protocol net.ProtocolNumber, data []byte, dst net.ProtocolAddress) error {
	// TODO
	return nil
}
