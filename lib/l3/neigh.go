package l3

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/vhqr0/gostack/lib/util"
)

type NeighTable struct {
	Mutex   sync.RWMutex
	Entries map[string]*NeighEntry
	// ip4 => Neigh4Key(ifidx, ip4) len=6
	// ip6 => Neigh6Key(ifidx, ip6) len=18
}

type NeighEntry struct {
	TS  int64
	MAC net.HardwareAddr
}

func (table *NeighTable) Update(key string, mac net.HardwareAddr) {
	table.Mutex.Lock()
	entry, ok := table.Entries[key]
	if !ok {
		entry = &NeighEntry{}
		table.Entries[key] = entry
	}
	entry.TS = time.Now().Unix()
	entry.MAC = mac
	table.Mutex.Unlock()
}

func (table *NeighTable) Lookup(key string) (mac net.HardwareAddr) {
	table.Mutex.RLock()
	entry, ok := table.Entries[key]
	if ok && entry != nil && entry.TS+60 > time.Now().Unix() {
		mac = entry.MAC
	}
	table.Mutex.RUnlock()
	return
}

func (stack *IPStack) Lookup(ifidx int, peer net.IP) net.HardwareAddr {
	key := NeighKey(ifidx, peer)

	if mac := stack.NeighTable.Lookup(key); mac != nil {
		return mac
	}

	switch len(peer) {
	case 4:
		stack.arpReqSend(ifidx, peer)
	case 16:
		stack.ndpNSSend(ifidx, peer)
	default:
		log.Panic(&util.InvalidIPLen{Len: len(peer)})
	}

	time.Sleep(10 * time.Millisecond)

	return stack.NeighTable.Lookup(key)
}

func NewNeighTable() *NeighTable {
	table := &NeighTable{
		Entries: make(map[string]*NeighEntry),
	}
	return table
}
