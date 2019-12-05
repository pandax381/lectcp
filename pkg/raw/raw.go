package raw

import (
	"io"
)

type Device interface {
	io.ReadWriteCloser
	Name() string
	Address() []byte
}
