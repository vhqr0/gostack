package udp

import (
	"sync"

	"github.com/vhqr0/gostack/lib/sock"
)

type UDPTable struct {
	Mutex   sync.RWMutex
	Entries map[string]*UDPSock
	// udp4 => UDP4Key(addr) len=6
	// udp6 => UDP6Key(addr) len=18
}

func (table *UDPTable) Lookup(addr *sock.Addr) *UDPSock {
	table.Mutex.RLock()
	sock := table.Entries[UDPKey(addr)]
	table.Mutex.RUnlock()
	return sock
}

func (table *UDPTable) Del(addr *sock.Addr) {
	table.Mutex.Lock()
	delete(table.Entries, UDPKey(addr))
	table.Mutex.Unlock()
}

func (table *UDPTable) Add(addr *sock.Addr, udpSock *UDPSock) (err error) {
	table.Mutex.Lock()
	key := UDPKey(addr)
	if _, ok := table.Entries[key]; ok {
		err = &sock.BusySockAddr{Addr: addr}
	} else {
		udpSock.Local = addr
		table.Entries[key] = udpSock
	}
	table.Mutex.Unlock()
	return
}

func NewUDPTable() *UDPTable {
	table := &UDPTable{
		Entries: make(map[string]*UDPSock),
	}
	return table
}
