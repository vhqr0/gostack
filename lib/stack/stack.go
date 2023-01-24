package stack

import (
	"github.com/vhqr0/gostack/lib/host"
	"github.com/vhqr0/gostack/lib/l2"
	"github.com/vhqr0/gostack/lib/l3"
	"github.com/vhqr0/gostack/lib/l4/udp"
	"github.com/vhqr0/gostack/lib/sock"
)

type Stack struct {
	Host *host.Host

	EthStack *l2.EthStack
	IPStack  *l3.IPStack
	UDPStack *udp.UDPStack

	sockFactory *sock.SockFactory
}

func (vstack *Stack) AutoRoute() {
	vstack.IPStack.AutoRoute()
}

func (vstack *Stack) AddRoute(ver int, ifname, peerStr, netStr string) error {
	return vstack.IPStack.AddRoute(ver, ifname, peerStr, netStr)
}

func (vstack *Stack) Run() {
	vstack.UDPStack.Run()
	vstack.IPStack.Run()
	vstack.EthStack.Run()
}

func (vstack *Stack) NewSock(family, typ uint32) (sock.Sock, error) {
	return vstack.sockFactory.NewSock(family, typ)
}

func NewStack(vhost *host.Host) *Stack {
	vstack := &Stack{
		Host: vhost,

		sockFactory: sock.NewSockFactory(),
	}

	vstack.EthStack = l2.NewEthStack(vhost)
	vstack.IPStack = l3.NewIPStack(vstack.EthStack)
	vstack.UDPStack = udp.NewUDPStack(vstack.IPStack, vstack.sockFactory)

	return vstack
}
