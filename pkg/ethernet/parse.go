package ethernet

import (
	"bytes"
	"encoding/binary"

	"github.com/pandax381/lectcp/pkg/net"
)

type header struct {
	Dst  Address
	Src  Address
	Type net.EthernetType
}

type frame struct {
	header
	payload []byte
}

func parse(data []byte) (*frame, error) {
	frame := frame{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &frame.header); err != nil {
		return nil, err
	}
	frame.payload = buf.Bytes()
	return &frame, nil
}
