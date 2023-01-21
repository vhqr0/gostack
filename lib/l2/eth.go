package l2

import (
	"net"

	"github.com/vhqr0/gostack/lib/host"
)

const (
	EthARP uint16 = 0x0806
	EthIP4 uint16 = 0x0800
	EthIP6 uint16 = 0x86dd
)

type EthStack struct {
	Verbose bool

	Host *host.Host

	MACSet map[string]struct{}

	recvChMap map[uint16]chan<- *EthPkt
	sendCh    chan *EthPkt
}

type EthPkt struct {
	IfIdx   int
	Peer    net.HardwareAddr
	Proto   uint16
	Payload []byte
}

func (stack *EthStack) Register(proto uint16, ch chan<- *EthPkt) {
	stack.recvChMap[proto] = ch
}

func (stack *EthStack) Send(pkt *EthPkt) {
	stack.sendCh <- pkt
}

func (stack *EthStack) Run() {
	go stack.ethSender()
	for ifidx := range stack.Host.Ifaces {
		go stack.ethReceiver(ifidx)
	}
}

func NewEthStack(vhost *host.Host) *EthStack {
	stack := &EthStack{
		Verbose: vhost.Verbose,

		Host: vhost,

		MACSet: NewMACSet(vhost),

		recvChMap: make(map[uint16]chan<- *EthPkt),
		sendCh:    make(chan *EthPkt, 1024),
	}
	return stack
}
