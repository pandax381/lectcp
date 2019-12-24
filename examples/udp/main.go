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

func init() {
	icmp.Init()
}

func setup() error {
	name := flag.String("name", "", "device name")
	addr := flag.String("addr", "", "hardware address")
	flag.Parse()
	raw, err := tuntap.NewTap(*name)
	if err != nil {
		return err
	}
	link, err := ethernet.NewDevice(raw)
	if err != nil {
		return err
	}
	if *addr != "" {
		link.SetAddress(ethernet.ParseAddress(*addr))
	}
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return err
	}
	iface, err := ip.CreateInterface(dev, "172.16.0.100", "255.255.255.0", "")
	if err != nil {
		panic(err)
	}
	dev.RegisterInterface(iface)
	return nil
}

func main() {
	if err := setup(); err != nil {
		panic(err)
	}
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
