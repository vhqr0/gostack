package tcp

import (
	"sync"

	"github.com/vhqr0/gostack/lib/sock"
)

type TCPTable struct {
	Mutex sync.RWMutex
	Entries map[string]*TCPSock
	// tcp4 listen => TCP4ListenKey(addr)      len=6
	// tcp4 conn   => TCP4ConnKey(local, peer) len=12
	// tcp6 listen => TCP6ListenKey(addr)      len=18
	// tcp6 conn   => TCP6ConnKey(local, peer) len=36
}

func (table *TCPTable) LookupListen(addr *sock.Addr) *TCPSock {
	table.Mutex.RLock()
	sock := table.Entries[TCPListenKey(addr)]
	table.Mutex.RUnlock()
	return sock
}

func (table *TCPTable) LookupConn(local, peer *sock.Addr) *TCPSock {
	table.Mutex.RLock()
	sock := table.Entries[TCPConnKey(local, peer)]
	table.Mutex.RUnlock()
	return sock
}

func (table *TCPTable) DelListen(addr *sock.Addr) {
	table.Mutex.Lock()
	delete(table.Entries, TCPListenKey(addr))
	table.Mutex.Unlock()
}

func (table *TCPTable) AddListen(addr *sock.Addr, tcpSock *TCPSock) (err error) {
	table.Mutex.Lock()
	key := TCPListenKey(addr)
	if _, ok := table.Entries[key]; ok {
		err = &sock.BusySockAddr{Addr: addr}
	} else {
		tcpSock.Local = addr
		table.Entries[key] = tcpSock
	}
	table.Mutex.Unlock()
	return
}

func (table *TCPTable) DelConn(local, peer *sock.Addr) {
	table.Mutex.Lock()
	delete(table.Entries, TCPConnKey(local, peer))
	table.Mutex.Unlock()
}

func (table *TCPTable) AddConn(local, peer *sock.Addr, tcpSock *TCPSock) (err error) {
	table.Mutex.Lock()
	key := TCPConnKey(local, peer)
	if _, ok := table.Entries[key]; ok {
		err = &sock.BusySockAddr{Addr: local}
	} else {
		tcpSock.Local = local
		table.Entries[key] = tcpSock
	}
	table.Mutex.Unlock()
	return
}

func NewTCPTable() *TCPTable {
	table := &TCPTable{
		Entries: make(map[string]*TCPSock),
	}
	return table
}
