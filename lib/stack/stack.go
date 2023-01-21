package stack

import (
	"github.com/vhqr0/gostack/lib/host"
	"github.com/vhqr0/gostack/lib/l2"
	"github.com/vhqr0/gostack/lib/l3"
)

type Stack struct {
	Host     *host.Host
	EthStack *l2.EthStack
	IPStack  *l3.IPStack
}

func (vstack *Stack) AutoRoute() {
	vstack.IPStack.AutoRoute()
}

func (vstack *Stack) AddRoute(ver int, ifname, peerStr, netStr string) error {
	return vstack.IPStack.AddRoute(ver, ifname, peerStr, netStr)
}

func (vstack *Stack) Run() {
	vstack.IPStack.Run()
	vstack.EthStack.Run()
}

func NewStack(vhost *host.Host) (vstack *Stack) {
	vstack = &Stack{Host: vhost}
	vstack.EthStack = l2.NewEthStack(vhost)
	vstack.IPStack = l3.NewIPStack(vstack.EthStack)
	return
}
