package main

import (
	"log"
	"time"

	"github.com/pandax381/lectcp/pkg/icmp"
	"github.com/pandax381/lectcp/pkg/ip"
	"github.com/pandax381/lectcp/pkg/loopback"
	"github.com/pandax381/lectcp/pkg/net"
)

func init() {
	icmp.Init()
}

func setup() error {
	link, err := loopback.NewDevice()
	if err != nil {
		return err
	}
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return err
	}
	log.Printf("[%s] %s\n", dev.Name(), dev.Address())
	iface, err := ip.CreateInterface(dev, "127.0.0.1", "255.0.0.0", "")
	if err != nil {
		return err
	}
	dev.RegisterInterface(iface)
	return nil
}

func main() {
	if err := setup(); err != nil {
		panic(err)
	}
	go func() {
		peer := ip.ParseAddress("127.0.0.1")
		for range time.Tick(3 * time.Second) {
			icmp.EchoRequest([]byte("1234567890"), peer)
		}
	}()
	for {
		time.Sleep(1 * time.Second)
	}
}
