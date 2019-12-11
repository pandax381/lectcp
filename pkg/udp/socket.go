package udp

import "fmt"

func Dial(local, remote *Address) (*Conn, error) {
	// TODO
	return nil, fmt.Errorf("dial failure")
}

func Listen(local *Address) (*Conn, error) {
	entry := repo.add(local)
	if entry == nil {
		return nil, fmt.Errorf("listen failure")
	}
	return &Conn{cb: entry}, nil
}
