package net

import (
	"fmt"
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
	running  bool
	mRunning sync.Mutex
	rxQueue  chan []byte
	term     chan chan struct{}
	ifaces   []ProtocolInterface
}

var devices = map[LinkDevice]*Device{}

func RegisterDevice(link LinkDevice) (*Device, error) {
	if _, exists := devices[link]; exists {
		return nil, fmt.Errorf("link device '%s' is already registered", link.Name())
	}
	dev := &Device{
		LinkDevice: link,
	}
	devices[link] = dev
	return dev, nil
}

func rxHandler(link LinkDevice, protocol EthernetType, payload []byte, src, dst HardwareAddress) {
	for k, entry := range protocols {
		if k == EthernetType(protocol) {
			log.Printf("rx: [%s] %s => %s (%s) %d bytes\n", link.Name(), src, dst, protocol, len(payload))
			entry.rxQueue <- &packet{devices[link.(LinkDevice)], payload, src, dst}
			return
		}
	}
}

func Devices() []*Device {
	ret := []*Device{}
	for _, v := range devices {
		ret = append(ret, v)
	}
	return ret
}

func (d *Device) Interfaces() []ProtocolInterface {
	return d.ifaces
}

func (d *Device) launchRxLoop() {
	go func() {
		var buf = make([]byte, d.HeaderSize()+d.MTU())
		for {
			n, err := d.Read(buf)
			if err != nil {
				close(d.rxQueue)
				return
			}
			d.rxQueue <- buf[:n]
		}
	}()
	for {
		select {
		case complete, ok := <-d.term:
			if ok {
				complete <- struct{}{}
			}
			return
		case buf, ok := <-d.rxQueue:
			if !ok {
				d.mRunning.Lock()
				d.running = false
				d.rxQueue = nil
				d.mRunning.Unlock()
				return
			}
			d.RxHandler(buf, rxHandler)
		}
	}
}

func (d *Device) Run() {
	d.mRunning.Lock()
	if d.running {
		d.mRunning.Unlock()
		return
	}
	d.running = true
	if d.rxQueue == nil {
		d.rxQueue = make(chan []byte)
	}
	d.mRunning.Unlock()
	go d.launchRxLoop()
}

func (d *Device) Stop() {
	d.mRunning.Lock()
	running := d.running
	d.mRunning.Unlock()
	if !running {
		return
	}
	complete := make(chan struct{})
	d.term <- complete
	<-complete
	d.mRunning.Lock()
	d.running = false
	d.mRunning.Unlock()
}

func (d *Device) RegisterInterface(iface ProtocolInterface) {
	d.ifaces = append(d.ifaces, iface)
}
