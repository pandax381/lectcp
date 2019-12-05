package ioctl

import (
	"syscall"
	"unsafe"
)

func SIOCGIFINDEX(name string) (int32, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [16]byte
		index int32
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFINDEX, uintptr(unsafe.Pointer(&ifreq)), 0, 0, 0); errno != 0 {
		return 0, errno
	}
	return ifreq.index, err
}

func SIOCGIFFLAGS(name string) (uint16, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifreq)), 0, 0, 0); errno != 0 {
		return 0, errno
	}
	return ifreq.flags, nil
}

func SIOCSIFFLAGS(name string, flags uint16) error {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	ifreq.flags = flags
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifreq)), 0, 0, 0); errno != 0 {
		return errno
	}
	return nil
}

type sockaddr struct {
	family uint16
	addr   [14]byte
}

func SIOCGIFHWADDR(name string) ([]byte, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name [syscall.IFNAMSIZ]byte
		addr sockaddr
		_pad [8]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFHWADDR, uintptr(unsafe.Pointer(&ifreq)), 0, 0, 0); errno != 0 {
		return nil, errno
	}
	return ifreq.addr.addr[:], nil
}
