package ip

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type Address [4]byte

var (
	EmptyAddress     = Address{0x00, 0x00, 0x00, 0x00}
	InvalidAddress   = Address{0x00, 0x00, 0x00, 0x00}
	BroadcastAddress = Address{0xff, 0xff, 0xff, 0xff}
)

func ParseAddress(s string) (Address, error) {
	parts := strings.FieldsFunc(s, func(c rune) bool {
		return c == '.'
	})
	if len(parts) != 4 {
		return InvalidAddress, fmt.Errorf("inconsistent parts: %s", s)
	}
	ret := Address{}
	for i, part := range parts {
		u, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return InvalidAddress, fmt.Errorf("invalid digits: %s", part)
		}
		ret[i] = uint8(u & 0xff)
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
	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
}

func (a Address) Uint32() uint32 {
	return *(*uint32)(unsafe.Pointer(&a[0]))
}

func (a Address) IsEmpty() bool {
	return a == EmptyAddress
}
