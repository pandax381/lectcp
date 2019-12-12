package udp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"unsafe"

	"github.com/pandax381/lectcp/pkg/ip"
	"github.com/pandax381/lectcp/pkg/net"
)

type Conn struct {
	cb   *cbEntry
	peer *Address
}

func (conn *Conn) Close() {
	// TODO
}

func (conn *Conn) Read(buf []byte) (int, error) {
	if conn.peer == nil {
		return -1, fmt.Errorf("this Conn is not dialed")
	}
	n, _, err := conn.ReadFrom(buf)
	return n, err
}

func (conn *Conn) ReadFrom(buf []byte) (int, *Address, error) {
	select {
	case q := <-conn.cb.rxQueue:
		n := copy(buf, q.data)
		return n, &Address{Addr: q.addr, Port: q.port}, nil
	}
}

func getAppropriateInterface(local, remote net.ProtocolAddress) net.ProtocolInterface {
	if local.IsEmpty() {
		return ip.GetInterfaceByRemoteAddress(remote)
	}
	return ip.GetInterface(local)
}

func (conn *Conn) Write(data []byte) error {
	if conn.peer == nil {
		return fmt.Errorf("this Conn is not dialed")
	}
	return conn.WriteTo(data, conn.peer)
}

func (conn *Conn) WriteTo(data []byte, peer *Address) error {
	hdr := header{}
	hdr.SourcePort = conn.cb.Port
	hdr.DestinationPort = peer.Port
	hdr.Length = uint16(int(unsafe.Sizeof(hdr)) + len(data))
	datagram := datagram{
		header: hdr,
		data:   data,
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &hdr)
	binary.Write(buf, binary.BigEndian, data)
	iface := getAppropriateInterface(conn.cb.Addr, peer.Addr)
	b := buf.Bytes()
	datagram.Checksum = net.Cksum16(b, len(b), pseudoHeaderSum(iface.Address(), peer.Addr, len(b)))
	binary.BigEndian.PutUint16(b[6:8], datagram.Checksum)
	log.Printf("[UDP] WriteTo: %s (%s:%d) => %s (%d bytes)\n", conn.cb.Address, iface.Address(), conn.cb.Address.Port, peer, len(data))
	datagram.dump()
	return iface.Tx(net.ProtocolNumberUDP, b, peer.Addr)
}
