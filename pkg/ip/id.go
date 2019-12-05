package ip

import "sync"

type idManager struct {
	num   uint16
	mutex sync.Mutex
}

func (id *idManager) next() uint16 {
	id.mutex.Lock()
	defer id.mutex.Unlock()
	id.num++
	return id.num
}

var idm = &idManager{}
