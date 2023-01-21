package l3

import (
	"net"

	"github.com/vhqr0/gostack/lib/l2"
)

const (
	IPICMP4 uint8 = 1
	IPICMP6 uint8 = 58
	IPTCP   uint8 = 6
	IPUDP   uint8 = 17
)

type IPStack struct {
	Verbose bool
	Forward bool

	EthStack *l2.EthStack

	IPSet       map[string]struct{}
	NeighTable  *NeighTable
	Route4Table *RouteTable
	Route6Table *RouteTable

	arpRecvCh   chan *l2.EthPkt
	ip4RecvCh   chan *l2.EthPkt
	ip6RecvCh   chan *l2.EthPkt
	icmp4RecvCh chan *IPPkt
	icmp6RecvCh chan *IPPkt

	recvChMap map[uint8]chan<- *IPPkt
	sendCh    chan *IPPkt
}

type IPPkt struct {
	IfIdx   int
	PeerMAC net.HardwareAddr
	Local   net.IP
	Peer    net.IP
	Proto   uint8
	Payload []byte
}

func (stack *IPStack) Register(proto uint8, ch chan<- *IPPkt) {
	stack.recvChMap[proto] = ch
}

func (stack *IPStack) Send(pkt *IPPkt) {
	stack.sendCh <- pkt
}

func (stack *IPStack) Run() {
	go stack.ipSender()
	go stack.arpReceiver()
	go stack.ip4Receiver()
	go stack.ip6Receiver()
	go stack.icmp4Receiver()
	go stack.icmp6Receiver()
}

func NewIPStack(ethStack *l2.EthStack) *IPStack {
	stack := &IPStack{
		Verbose: ethStack.Verbose,
		Forward: ethStack.Host.Forward,

		EthStack: ethStack,

		IPSet:       NewIPSet(ethStack.Host),
		NeighTable:  NewNeighTable(),
		Route4Table: &RouteTable{},
		Route6Table: &RouteTable{},

		arpRecvCh:   make(chan *l2.EthPkt, 1024),
		ip4RecvCh:   make(chan *l2.EthPkt, 1024),
		ip6RecvCh:   make(chan *l2.EthPkt, 1024),
		icmp4RecvCh: make(chan *IPPkt, 1024),
		icmp6RecvCh: make(chan *IPPkt, 1024),

		recvChMap: make(map[uint8]chan<- *IPPkt),
		sendCh:    make(chan *IPPkt, 1024),
	}

	stack.Register(IPICMP4, stack.icmp4RecvCh)
	stack.Register(IPICMP6, stack.icmp6RecvCh)

	ethStack.Register(l2.EthARP, stack.arpRecvCh)
	ethStack.Register(l2.EthIP4, stack.ip4RecvCh)
	ethStack.Register(l2.EthIP6, stack.ip6RecvCh)

	return stack
}
