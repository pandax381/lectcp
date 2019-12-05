package net

type HardwareType uint16

const (
	HardwareTypeEthernet = 0x0001
)

func (t HardwareType) String() string {
	switch t {
	case HardwareTypeEthernet:
		return "Ethernet"
	default:
		return "Unknown"
	}
}

type EthernetType uint16

const (
	EthernetTypeIP   EthernetType = 0x0800
	EthernetTypeARP  EthernetType = 0x0806
	EthernetTypeIPv6 EthernetType = 0x86dd
)

func (t EthernetType) String() string {
	switch t {
	case EthernetTypeIP:
		return "IP"
	case EthernetTypeARP:
		return "ARP"
	case EthernetTypeIPv6:
		return "IPv6"
	default:
		return "Unknown"
	}
}
