package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/pandax381/lectcp/pkg/arp"
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

func (iface *Interface) xmit(datagram *datagram, nexthop net.ProtocolAddress) error {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &datagram.header)
	binary.Write(buf, binary.BigEndian, datagram.payload)
	b := buf.Bytes()
	binary.BigEndian.PutUint16(b[10:12], net.Cksum16(b, len(b), 0))
	var hardwareAddress []byte
	if iface.Device().NeedARP() {
		if nexthop != nil {
			var err error
			hardwareAddress, err = arp.Resolve(iface, nexthop.Bytes(), b)
			if err != nil {
				return err
			}
		} else {
			hardwareAddress = iface.Device().BroadcastAddress().Bytes()
		}
	}
	return iface.Device().Tx(net.EthernetTypeIP, b, hardwareAddress)
}

func (iface *Interface) Tx(protocol net.ProtocolNumber, data []byte, dst net.ProtocolAddress) error {
	log.Printf("tx: [%s] %s => %s (%s) %d bytes\n", iface.Device().Name(), iface.Address(), dst, protocol, len(data))
	var nexthop net.ProtocolAddress
	src := iface.unicast
	if dst.(Address) != BroadcastAddress {
		network := Address{
			iface.unicast[0] & iface.netmask[0],
			iface.unicast[1] & iface.netmask[1],
			iface.unicast[2] & iface.netmask[2],
			iface.unicast[3] & iface.netmask[3],
		}
		if dst.(Address).Uint32()&iface.netmask.Uint32() != network.Uint32() {
			if iface.gateway == EmptyAddress {
				return fmt.Errorf("no route to host")
			}
			nexthop = iface.gateway
		}
	}
	assembler := newAssembler(protocol, data, src, dst, idm.next(), iface.Device().MTU())
	for _, datagram := range assembler.assemble() {
		if err := iface.xmit(datagram, nexthop); err != nil {
			return err
		}
	}
	return nil
}
