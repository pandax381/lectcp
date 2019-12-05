package main

import (
	"log"
	"time"

	"github.com/pandax381/lectcp/pkg/icmp"
	"github.com/pandax381/lectcp/pkg/ip"
	"github.com/pandax381/lectcp/pkg/loopback"
	"github.com/pandax381/lectcp/pkg/net"
)

func setupLoopback() (*net.Device, error) {
	link, err := loopback.NewDevice()
	if err != nil {
		return nil, err
	}
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return nil, err
	}
	return dev, nil
}

func init() {
	icmp.Init()
}

func main() {
	dev, err := setupLoopback()
	if err != nil {
		panic(err)
	}
	log.Printf("[%s] %s\n", dev.Name(), dev.Address())
	iface, err := ip.CreateInterface(dev, "127.0.0.1", "255.0.0.0", "")
	if err != nil {
		panic(err)
	}
	dev.RegisterInterface(iface)
	dev.Run()
	addr, err := ip.ParseAddress("127.0.0.1")
	if err != nil {
		panic(err)
	}
	for {
		icmp.EchoRequest(iface, []byte("1234567890"), addr)
		time.Sleep(1 * time.Second)
	}
}
