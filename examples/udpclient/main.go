package main

import (
	"flag"
	"time"

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

func setupTap() error {
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
	if err := setupTap(); err != nil {
		panic(err)
	}
	conn, err := udp.Dial(
		nil,
		&udp.Address{
			Addr: ip.ParseAddress("172.16.0.1"),
			Port: 10381,
		},
	)
	if err != nil {
		panic(err)
	}
	for {
		conn.Write([]byte("hoge\n"))
		time.Sleep(3 * time.Second)
	}
}
