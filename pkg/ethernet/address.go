package ethernet

import (
	"fmt"
	"strconv"
	"strings"
)

type Address [6]byte

var (
	EmptyAddress     = Address{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	InvalidAddress   = Address{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	BroadcastAddress = Address{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

func ParseAddress(s string) (Address, error) {
	parts := strings.FieldsFunc(s, func(c rune) bool {
		return c == ':' || c == '-'
	})
	ret := Address{}
	if len(parts) != 6 {
		return ret, fmt.Errorf("inconsistent parts: %s", s)
	}
	for i, part := range parts {
		u, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return ret, fmt.Errorf("invalid hex digits: %s", s)
		}
		ret[i] = byte(u)
	}
	return ret, nil
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) Len() uint8 {
	return uint8(len(a))
}

func (a Address) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", a[0], a[1], a[2], a[3], a[4], a[5])
}
