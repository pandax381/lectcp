package main

import (
	"flag"
	"time"

	"github.com/pandax381/lectcp/pkg/ethernet"
	"github.com/pandax381/lectcp/pkg/net"
	"github.com/pandax381/lectcp/pkg/raw/tuntap"
)

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
	net.RegisterDevice(link)
	return nil
}

func main() {
	if err := setup(); err != nil {
		panic(err)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
