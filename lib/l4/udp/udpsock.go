package udp

import (
	"sync"

	"github.com/vhqr0/gostack/lib/sock"
)

const (
	UDPClosed = iota
	UDPOpen
)

type UDPSock struct {
	Stack *UDPStack

	Family uint32
	Typ    uint32

	Mutex  sync.RWMutex
	Status int
	Local  *sock.Addr
	Peer   *sock.Addr

	recvCh chan *UDPPkt
}

func (udpSock *UDPSock) ValidateAddr(addr *sock.Addr) error {
	if addr == nil || addr.Port == 0 {
		return &sock.InvalidSockAddr{Addr: addr}
	}
	switch udpSock.Family {
	case sock.AFIP4:
		if len(addr.IP) != 4 {
			return &sock.InvalidSockAddr{Addr: addr}
		}
	case sock.AFIP6:
		if len(addr.IP) != 16 {
			return &sock.InvalidSockAddr{Addr: addr}
		}
	default:
		return &sock.InvalidSockFamilyOrTyp{Family: udpSock.Family, Typ: udpSock.Typ}
	}
	if !addr.IP.IsLoopback() && !addr.IP.IsGlobalUnicast() {
		return &sock.InvalidSockAddr{Addr: addr}
	}
	return nil
}

func (stack *UDPStack) NewUDPSock(family, typ uint32) (*UDPSock, error) {
	sock := &UDPSock{
		Stack: stack,

		Family: family,
		Typ:    typ,

		Status: UDPOpen,

		recvCh: make(chan *UDPPkt), // Notice: blocking chan
	}
	return sock, nil
}

func (stack *UDPStack) NewSock(family, typ uint32) (sock.Sock, error) {
	return stack.NewUDPSock(family, typ)
}
