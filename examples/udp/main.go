package main

import (
	"flag"
	"log"

	"github.com/pandax381/lectcp/pkg/ethernet"
	"github.com/pandax381/lectcp/pkg/icmp"
	"github.com/pandax381/lectcp/pkg/ip"
	"github.com/pandax381/lectcp/pkg/net"
	"github.com/pandax381/lectcp/pkg/raw/tuntap"
	"github.com/pandax381/lectcp/pkg/udp"
)

func setupTap(name, hwaddr string) (*net.Device, error) {
	raw, err := tuntap.NewTap(name)
	if err != nil {
		return nil, err
	}
	link, err := ethernet.NewDevice(raw)
	if err != nil {
		return nil, err
	}
	if hwaddr != "" {
		addr, err := ethernet.ParseAddress(hwaddr)
		if err != nil {
			return nil, err
		}
		link.SetAddress(addr)
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
	name := flag.String("name", "", "device name")
	addr := flag.String("addr", "", "hardware address")
	flag.Parse()
	dev, err := setupTap(*name, *addr)
	if err != nil {
		panic(err)
	}
	log.Printf("[%s] %s\n", dev.Name(), dev.Address())
	iface, err := ip.CreateInterface(dev, "172.16.0.100", "255.255.255.0", "")
	if err != nil {
		panic(err)
	}
	dev.RegisterInterface(iface)
	dev.Run()
	conn, err := udp.Listen(
		&udp.Address{
			Addr: ip.EmptyAddress,
			Port: 7,
		},
	)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 65536)
	for {
		n, peer, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		log.Printf("main: receive %d bytes data from %s", n, peer)
		conn.WriteTo(buf[:n], peer)
	}
}
