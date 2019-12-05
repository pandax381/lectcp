package main

import (
	"flag"
	"log"
	"time"

	"github.com/pandax381/lectcp/pkg/arp"
	"github.com/pandax381/lectcp/pkg/ethernet"
	"github.com/pandax381/lectcp/pkg/ip"
	"github.com/pandax381/lectcp/pkg/net"
	"github.com/pandax381/lectcp/pkg/raw/pfpacket"
	"github.com/pandax381/lectcp/pkg/raw/tuntap"
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

func setupPFPacket(name, hwaddr string) (*net.Device, error) {
	raw, err := pfpacket.NewPFPacket(name)
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
	arp.Init()
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
	for {
		time.Sleep(1 * time.Second)
	}
}
