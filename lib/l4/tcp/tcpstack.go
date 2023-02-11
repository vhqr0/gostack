package tcp

import (
	"github.com/vhqr0/gostack/lib/l3"
	"github.com/vhqr0/gostack/lib/sock"
)

type TCPStack struct {
	Verbose bool

	IPStack *l3.IPStack

	SockTable *TCPTable

	recvCh chan *l3.IPPkt
}

func (stack *TCPStack) Run() {
	go stack.tcpReceiver()
}

func NewTCPStack(ipStack *l3.IPStack, sockFactory *sock.SockFactory) *TCPStack {
	stack := &TCPStack{
		Verbose: ipStack.Verbose,

		IPStack: ipStack,

		SockTable: NewTCPTable(),

		recvCh: make(chan *l3.IPPkt, 1024),
	}

	ipStack.Register(l3.IPTCP, stack.recvCh)

	sockFactory.Register(sock.AFIP4, sock.SockStream, stack)
	sockFactory.Register(sock.AFIP6, sock.SockStream, stack)

	return stack
}
