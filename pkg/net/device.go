package net

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type LinkDeviceCallbackHandler func(link LinkDevice, protocol EthernetType, payload []byte, src, dst HardwareAddress)

type LinkDevice interface {
	Type() HardwareType
	Name() string
	Address() HardwareAddress
	BroadcastAddress() HardwareAddress
	MTU() int
	HeaderSize() int
	NeedARP() bool
	Close()
	Read(data []byte) (int, error)
	RxHandler(frame []byte, callback LinkDeviceCallbackHandler)
	Tx(proto EthernetType, data []byte, dst []byte) error
}

type Device struct {
	LinkDevice
	errors chan error
	ifaces []ProtocolInterface
	sync.RWMutex
}

var devices = sync.Map{}

func RegisterDevice(link LinkDevice) (*Device, error) {
	if _, exists := devices.Load(link); exists {
		return nil, fmt.Errorf("link device '%s' is already registered", link.Name())
	}
	dev := &Device{
		LinkDevice: link,
		errors:     make(chan error),
	}
	// launch rx loop
	go func() {
		var buf = make([]byte, dev.HeaderSize()+dev.MTU())
		for {
			n, err := dev.Read(buf)
			if n > 0 {
				dev.RxHandler(buf[:n], rxHandler)
			}
			if err != nil {
				dev.errors <- err
				break
			}
		}
		close(dev.errors)
	}()
	devices.Store(link, dev)
	return dev, nil
}

func rxHandler(link LinkDevice, protocol EthernetType, payload []byte, src, dst HardwareAddress) {
	protocols.Range(func(k interface{}, v interface{}) bool {
		var (
			Type  = k.(EthernetType)
			entry = v.(*entry)
		)
		if Type == EthernetType(protocol) {
			dev, ok := devices.Load(link)
			if !ok {
				panic("device not found")
			}
			entry.rxQueue <- &packet{
				dev:  dev.(*Device),
				data: payload,
				src:  src,
				dst:  dst,
			}
			return false // break range loop
		}
		return true
	})
}

func Devices() []*Device {
	ret := []*Device{}
	devices.Range(func(_, v interface{}) bool {
		ret = append(ret, v.(*Device))
		return true
	})
	return ret
}

func (d *Device) Interfaces() []ProtocolInterface {
	d.RLock()
	ret := make([]ProtocolInterface, len(d.ifaces))
	for i, iface := range d.ifaces {
		ret[i] = iface
	}
	d.RUnlock()
	return ret
}

func (d *Device) Shutdown() {
	d.LinkDevice.Close()
	if err := <-d.errors; err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}
	devices.Delete(d.LinkDevice)
}

func (d *Device) RegisterInterface(iface ProtocolInterface) {
	d.Lock()
	d.ifaces = append(d.ifaces, iface)
	d.Unlock()
}
