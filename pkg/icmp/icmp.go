package icmp

import (
	"log"

	"github.com/pandax381/lectcp/pkg/ip"
	"github.com/pandax381/lectcp/pkg/net"
)

func init() {
	ip.RegisterProtocol(net.ProtocolNumberICMP, rxHandler)
}

func Init() {
	// do nothing
}

func rxHandler(iface net.ProtocolInterface, data []byte, src, dst net.ProtocolAddress) error {
	msg, err := parse(data)
	if err != nil {
		return err
	}
	log.Printf("rx: %s => %s (%s)\n", src, dst, msg.messageType())
	msg.dump()
	switch msg.messageType() {
	case messageTypeEcho:
		request := msg.(*messageEcho)
		reply := &messageEcho{
			header: header{messageTypeEchoReply, 0, 0},
			id:     request.id,
			seq:    request.seq,
			data:   request.data,
		}
		return tx(iface, reply, src)
	case messageTypeEchoReply:
		// do nothing
	}
	return nil
}

func tx(iface net.ProtocolInterface, msg message, dst net.ProtocolAddress) error {
	buf, err := msg.marshal()
	if err != nil {
		return err
	}
	log.Printf("tx: %s => %s (%s)\n", iface.Address(), dst, msg.messageType())
	msg.dump()
	return iface.Tx(net.ProtocolNumberICMP, buf, dst)
}

func EchoRequest(data []byte, dst net.ProtocolAddress) error {
	request := &messageEcho{
		header: header{messageTypeEcho, 0, 0},
		id:     0,
		seq:    0,
		data:   data,
	}
	return tx(ip.GetInterfaceByRemoteAddress(dst), request, dst)
}
