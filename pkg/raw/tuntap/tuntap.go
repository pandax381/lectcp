package tuntap

import (
	"io"
)

type Tap struct {
	io.ReadWriteCloser
	name string
}

func NewTap(name string) (*Tap, error) {
	n, f, err := openTap(name)
	if err != nil {
		return nil, err
	}
	return &Tap{
		ReadWriteCloser: f,
		name:            n,
	}, nil
}

func (t Tap) Address() []byte {
	addr, _ := getAddress(t.name)
	return addr[:6]
}

func (t Tap) Name() string {
	return t.name
}
