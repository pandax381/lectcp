package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pandax381/lectcp/pkg/ethernet"
	"github.com/pandax381/lectcp/pkg/icmp"
	"github.com/pandax381/lectcp/pkg/ip"
	"github.com/pandax381/lectcp/pkg/net"
	"github.com/pandax381/lectcp/pkg/raw/tuntap"
	"github.com/pandax381/lectcp/pkg/udp"
)

var sig chan os.Signal

func init() {
	icmp.Init()
}

func setup() error {
	// signal handling
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// parse command line params
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
		return err
	}
	dev.RegisterInterface(iface)
	return nil
}

func main() {
	err := setup()
	if err != nil {
		panic(err)
	}
	peer := &udp.Address{
		Addr: ip.ParseAddress("172.16.0.1"),
		Port: 10381,
	}
	conn, err := udp.Dial(nil, peer)
	if err != nil {
		panic(err)
	}
	// launch send loop
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		data := []byte("hoge\n")
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				fmt.Printf("%d bytes data send to %s\n", len(data), peer)
				conn.Write(data)
			}
		}
	}()
	select {
	case s := <-sig:
		fmt.Printf("sig: %s\n", s)
		cancel()
	}
	wg.Wait()
	conn.Close()
	fmt.Println("good bye")
}
