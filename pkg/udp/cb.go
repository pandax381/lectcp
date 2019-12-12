package udp

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/pandax381/lectcp/pkg/net"
)

type queueEntry struct {
	addr net.ProtocolAddress
	port uint16
	data []byte
}

type cbEntry struct {
	*Address
	rxQueue chan *queueEntry
}

type cbRepository struct {
	list  *list.List
	mutex sync.RWMutex
}

var repo *cbRepository

func newCbRepository() *cbRepository {
	return &cbRepository{
		list: list.New(),
	}
}

func (repo *cbRepository) lookupUnlocked(addr *Address) *cbEntry {
	for elem := repo.list.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*cbEntry)
		if entry.Port == addr.Port && (entry.Addr.IsEmpty() || entry.Addr == addr.Addr) {
			return entry
		}
	}
	return nil
}

func (repo *cbRepository) lookup(addr *Address) *cbEntry {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	return repo.lookupUnlocked(addr)
}

func (repo *cbRepository) getAvailablePort(addr net.ProtocolAddress) uint16 {
	var port uint16
	for port = 40000; port <= 65535; port++ {
		var elem *list.Element
		for elem = repo.list.Front(); elem != nil; elem = elem.Next() {
			entry := elem.Value.(*cbEntry)
			if entry.Port == port && (entry.Addr.IsEmpty() || entry.Addr == addr) {
				break
			}
		}
		if elem == nil {
			return port
		}
	}
	return 0
}

func (repo *cbRepository) add(addr *Address) *cbEntry {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	if addr.Port == 0 {
		addr.Port = repo.getAvailablePort(addr.Addr)
		if addr.Port == 0 {
			return nil
		}
	} else {
		if repo.lookupUnlocked(addr) != nil {
			fmt.Println("entry exists")
			return nil
		}
	}
	entry := &cbEntry{
		Address: addr,
		rxQueue: make(chan *queueEntry),
	}
	repo.list.PushBack(entry)
	return entry
}

func (repo *cbRepository) del(entry *cbEntry) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	for elem := repo.list.Front(); elem != nil; elem = elem.Next() {
		if elem.Value == entry {
			repo.list.Remove(elem)
			return true
		}
	}
	return false
}
