package net

type HardwareAddress interface {
	Bytes() []byte
	Len() uint8
	String() string
}

type ProtocolAddress interface {
	Bytes() []byte
	Len() uint8
	String() string
	IsEmpty() bool
}
