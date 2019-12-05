package arp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/pandax381/lectcp/pkg/net"
)

const (
	operationRequest = 1
	operationReply   = 2
)

func init() {
	repo = newArpTable()
	net.RegisterProtocol(net.EthernetTypeARP, rxHandler)
}

func Init() {
	// do nothing
}

func rxHandler(dev *net.Device, data []byte, src, dst net.HardwareAddress) error {
	msg, err := parse(data)
	if err != nil {
		return err
	}
	log.Printf("%s => %s (%d bytes)\n", src, dst, len(data))
	marge := repo.update(msg.sourceProtocolAddress, msg.sourceHardwareAddress)
	for _, iface := range dev.Interfaces() {
		if bytes.Compare(msg.targetProtocolAddress, iface.Address().Bytes()) == 0 {
			if !marge {
				repo.insert(iface, msg.sourceProtocolAddress, msg.sourceHardwareAddress)
			}
			if msg.OperationCode == operationRequest {
				if err = reply(iface, msg.sourceProtocolAddress, msg.sourceHardwareAddress); err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}

func reply(iface net.ProtocolInterface, targetProtocolAddress []byte, targetHardwareAddress []byte) error {
	dev := iface.Device()
	hdr := header{
		HardwareType:          dev.Type(),
		ProtocolType:          iface.Type(),
		HardwareAddressLength: dev.Address().Len(),
		ProtocolAddressLength: iface.Address().Len(),
		OperationCode:         operationReply,
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &hdr)
	binary.Write(buf, binary.BigEndian, dev.Address().Bytes())
	binary.Write(buf, binary.BigEndian, iface.Address().Bytes())
	binary.Write(buf, binary.BigEndian, targetHardwareAddress)
	binary.Write(buf, binary.BigEndian, targetProtocolAddress)

	return dev.Tx(net.EthernetTypeARP, buf.Bytes(), targetHardwareAddress)
}

func Resolve(iface net.ProtocolInterface, target []byte, data []byte) ([]byte, error) {
	repo.mutex.RLock()
	repo.mutex.RUnlock()
	entry := repo.lookupUnlocked(target)
	if entry == nil {
		// TODO: send arp request
		return nil, fmt.Errorf("arp entry not found")
	}
	return entry.hardwareAddress, nil
}
