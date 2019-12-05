package main

import (
	"flag"
	"log"
	"time"

	"github.com/pandax381/lectcp/pkg/ethernet"
	"github.com/pandax381/lectcp/pkg/net"
	"github.com/pandax381/lectcp/pkg/raw/pfpacket"
	"github.com/pandax381/lectcp/pkg/raw/tuntap"
)

func setupTap(name, hwaddr string) error {
	raw, err := tuntap.NewTap(name)
	if err != nil {
		return err
	}
	link, err := ethernet.NewDevice(raw)
	if err != nil {
		return err
	}
	if hwaddr != "" {
		addr, err := ethernet.ParseAddress(hwaddr)
		if err != nil {
			return err
		}
		link.SetAddress(addr)
	}
	log.Printf("%s [%s]\n", link.Name(), link.Address())
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return err
	}
	dev.Run()
	return nil
}

func setupPFPacket(name, hwaddr string) error {
	raw, err := pfpacket.NewPFPacket(name)
	if err != nil {
		return err
	}
	link, err := ethernet.NewDevice(raw)
	if err != nil {
		return err
	}
	if hwaddr != "" {
		addr, err := ethernet.ParseAddress(hwaddr)
		if err != nil {
			return err
		}
		link.SetAddress(addr)
	}
	log.Printf("[%s] %s\n", link.Name(), link.Address())
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return err
	}
	dev.Run()
	return nil
}

func main() {
	name := flag.String("name", "", "device name")
	addr := flag.String("addr", "", "hardware address")
	flag.Parse()
	if err := setupTap(*name, *addr); err != nil {
		panic(err)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
