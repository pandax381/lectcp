package tuntap

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/pandax381/lectcp/pkg/ioctl"
)

const cloneDevice = "/dev/net/tun"

func openTap(name string) (string, *os.File, error) {
	if len(name) >= syscall.IFNAMSIZ {
		return "", nil, fmt.Errorf("name is too long")
	}
	file, err := os.OpenFile(cloneDevice, os.O_RDWR, 0600)
	if err != nil {
		return "", nil, err
	}
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], []byte(name))
	ifreq.flags = syscall.IFF_TAP | syscall.IFF_NO_PI
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(file.Fd()), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifreq)), 0, 0, 0); errno != 0 {
		file.Close()
		return "", nil, errno
	}
	name = string(ifreq.name[:bytes.IndexByte(ifreq.name[:], 0)])
	flags, err := ioctl.SIOCGIFFLAGS(name)
	if err != nil {
		file.Close()
		return "", nil, err
	}
	flags |= (syscall.IFF_UP | syscall.IFF_RUNNING)
	if err := ioctl.SIOCSIFFLAGS(name, flags); err != nil {
		file.Close()
		return "", nil, err
	}
	return name, file, nil
}

func getAddress(name string) ([]byte, error) {
	addr, err := ioctl.SIOCGIFHWADDR(name)
	if err != nil {
		return nil, err
	}
	return addr, nil
}
