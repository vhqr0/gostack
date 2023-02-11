package udp

import (
	"github.com/vhqr0/gostack/lib/l3"
	"github.com/vhqr0/gostack/lib/sock"
)

type UDPStack struct {
	Verbose bool

	IPStack *l3.IPStack

	SockTable *UDPTable

	recvCh chan *l3.IPPkt
}

func (stack *UDPStack) Run() {
	go stack.udpReceiver()
}

func NewUDPStack(ipStack *l3.IPStack, sockFactory *sock.SockFactory) *UDPStack {
	stack := &UDPStack{
		Verbose: ipStack.Verbose,

		IPStack: ipStack,

		SockTable: NewUDPTable(),

		recvCh: make(chan *l3.IPPkt, 1024),
	}

	ipStack.Register(l3.IPUDP, stack.recvCh)

	sockFactory.Register(sock.AFIP4, sock.SockDgram, stack)
	sockFactory.Register(sock.AFIP6, sock.SockDgram, stack)

	return stack
}
