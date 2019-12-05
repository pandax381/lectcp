package arp

import (
	"bytes"
	"sync"
	"time"

	"github.com/pandax381/lectcp/pkg/net"
)

type arpEntry struct {
	protocolAddress []byte
	hardwareAddress []byte
	iface           net.ProtocolInterface
	timestamp       time.Time
}

type arpTable struct {
	storage []*arpEntry
	mutex   sync.RWMutex
}

var repo *arpTable

func newArpTable() *arpTable {
	return &arpTable{
		storage: make([]*arpEntry, 0, 1024),
	}
}

func (tbl *arpTable) lookupUnlocked(protocolAddress []byte) *arpEntry {
	for _, entry := range tbl.storage {
		if bytes.Compare(entry.protocolAddress, protocolAddress) == 0 {
			return entry
		}
	}
	return nil
}

func (tbl *arpTable) lookup(protocolAddress []byte) *arpEntry {
	tbl.mutex.RLock()
	defer tbl.mutex.RUnlock()
	return tbl.lookupUnlocked(protocolAddress)
}

func (tbl *arpTable) update(protocolAddress []byte, hardwareAddress []byte) bool {
	tbl.mutex.Lock()
	defer tbl.mutex.Unlock()
	entry := tbl.lookupUnlocked(protocolAddress)
	if entry == nil {
		return false
	}
	entry.hardwareAddress = hardwareAddress
	entry.timestamp = time.Now()
	return true
}

func (tbl *arpTable) insert(iface net.ProtocolInterface, protocolAddress []byte, hardwareAddress []byte) bool {
	tbl.mutex.Lock()
	defer tbl.mutex.Unlock()
	if tbl.lookupUnlocked(protocolAddress) != nil {
		return false
	}
	entry := &arpEntry{
		protocolAddress: protocolAddress,
		hardwareAddress: hardwareAddress,
		iface:           iface,
		timestamp:       time.Now(),
	}
	tbl.storage = append(tbl.storage, entry)
	return true
}

func (tbl *arpTable) length() int {
	tbl.mutex.RLock()
	defer tbl.mutex.RUnlock()
	return len(tbl.storage)
}
