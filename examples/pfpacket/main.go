package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/pandax381/lectcp/pkg/raw/pfpacket"
)

func main() {
	name := flag.String("name", "", "device name")
	flag.Parse()
	dev, err := pfpacket.NewPFPacket(*name)
	if err != nil {
		panic(err)
	}
	fmt.Printf("name=%s, addr=%x\n", dev.Name(), dev.Address())
	buf := make([]byte, 4096)
	for {
		n, err := dev.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("--- [%s] incomming %d bytes data ---\n", dev.Name(), n)
		fmt.Printf("%s", hex.Dump(buf[:n]))
	}
}
