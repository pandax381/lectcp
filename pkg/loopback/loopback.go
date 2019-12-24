package loopback

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/pandax381/lectcp/pkg/net"
)

type header struct {
	Type net.EthernetType
}

type Device struct {
	name  string
	mtu   int
	queue chan []byte
}

var _ net.LinkDevice = &Device{} // interface check

var dev = Device{
	name:  "loopback0",
	mtu:   65536,
	queue: make(chan []byte),
}

func NewDevice() (*Device, error) {
	return &dev, nil
}

func (d *Device) Type() net.HardwareType {
	return net.HardwareTypeLoopback
}

func (d *Device) Name() string {
	return d.name
}

func (d *Device) Address() net.HardwareAddress {
	return nil
}

func (d *Device) BroadcastAddress() net.HardwareAddress {
	return nil
}

func (d *Device) MTU() int {
	return d.mtu
}

func (d *Device) HeaderSize() int {
	return int(unsafe.Sizeof(header{}))
}

func (d *Device) NeedARP() bool {
	return false
}

func (d *Device) Close() {
	close(d.queue)
}

func (d *Device) Read(buf []byte) (int, error) {
	var err error
	data, ok := <-d.queue
	if !ok {
		err = io.EOF
	}
	return copy(buf, data), err
}

func (d *Device) RxHandler(data []byte, callback net.LinkDeviceCallbackHandler) {
	hdr := header{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return
	}
	callback(d, hdr.Type, buf.Bytes(), nil, nil)
}

func (d *Device) Tx(Type net.EthernetType, data []byte, dst []byte) error {
	buf := make([]byte, 2+len(data))
	binary.BigEndian.PutUint16(buf[0:2], uint16(Type))
	copy(buf[2:], data)
	d.queue <- buf
	return nil
}
