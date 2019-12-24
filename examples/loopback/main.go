package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pandax381/lectcp/pkg/arp"
	"github.com/pandax381/lectcp/pkg/icmp"
	"github.com/pandax381/lectcp/pkg/ip"
	"github.com/pandax381/lectcp/pkg/loopback"
	"github.com/pandax381/lectcp/pkg/net"
)

var sig chan os.Signal

func init() {
	arp.Init()
	icmp.Init()
}

func setup() (*net.Device, error) {
	// signal handling
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	link, err := loopback.NewDevice()
	if err != nil {
		return nil, err
	}
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return nil, err
	}
	iface, err := ip.CreateInterface(dev, "127.0.0.1", "255.0.0.0", "")
	if err != nil {
		return nil, err
	}
	dev.RegisterInterface(iface)
	return dev, nil
}

func main() {
	dev, err := setup()
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] %s\n", dev.Name(), dev.Address())
	for _, iface := range dev.Interfaces() {
		fmt.Printf("  - %s\n", iface.Address())
	}
	// launch send loop
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		peer := ip.ParseAddress("127.0.0.1")
		data := []byte("1234567890")
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				fmt.Printf("send ICMP Echo to %s\n", peer)
				icmp.EchoRequest(data, peer)
			}
		}
	}()
	select {
	case s := <-sig:
		fmt.Printf("sig: %s\n", s)
		cancel()
	}
	wg.Wait()
	dev.Shutdown()
	fmt.Println("good bye")
}
