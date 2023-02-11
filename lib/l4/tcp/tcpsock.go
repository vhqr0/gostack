package tcp

import (
	"sync"

	"github.com/vhqr0/gostack/lib/sock"
)

const (
	TCPOpen = iota
	TCPClosed
	TCPEstablished

	TCPListen
	TCPSynSent
	TCPSynRcvd

	TCPFinWait1
	TCPFinWait2
	TCPClosing
	TCPTimeWait

	TCPCloseWait
	TCPLastAck
)

type TCPSock struct {
	Stack *TCPStack

	Family uint32
	Typ    uint32

	Mutex  sync.RWMutex
	Status int
	Local  *sock.Addr
	Peer   *sock.Addr
}

func (stack *TCPStack) NewTCPSock(family, typ uint32) (*TCPSock, error) {
	sock := &TCPSock{
		Stack: stack,

		Family: family,
		Typ:    typ,

		Status: TCPOpen,
	}

	return sock, nil
}

func (tcpSock *TCPSock) ValidateAddr(addr *sock.Addr) error {
	if addr == nil || addr.Port == 0 {
		return &sock.InvalidSockAddr{Addr: addr}
	}
	switch tcpSock.Family {
	case sock.AFIP4:
		if len(addr.IP) != 4 {
			return &sock.InvalidSockAddr{Addr: addr}
		}
	case sock.AFIP6:
		if len(addr.IP) != 16 {
			return &sock.InvalidSockAddr{Addr: addr}
		}
	default:
		return &sock.InvalidSockArgs{Family: tcpSock.Family, Typ: tcpSock.Typ}
	}
	if !addr.IP.IsLoopback() && !addr.IP.IsGlobalUnicast() {
		return &sock.InvalidSockAddr{Addr: addr}
	}
	return nil
}

func (stack *TCPStack) NewSock(family, typ uint32) (sock.Sock, error) {
	return stack.NewTCPSock(family, typ)
}
