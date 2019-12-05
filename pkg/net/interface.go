package net

type ProtocolInterface interface {
	Type() EthernetType
	Address() ProtocolAddress
	Tx(protocol ProtocolNumber, data []byte, dst ProtocolAddress) error
	Device() *Device
}
