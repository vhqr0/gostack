package l3

import (
	"encoding/binary"
	"log"

	"github.com/vhqr0/gostack/lib/util"
)

const (
	ICMP4EchoReq uint8 = 8
	ICMP4EchoRep uint8 = 0
)

func (stack *IPStack) icmp4Receiver() {
	// typ   uint8
	// code  uint8
	// cksum uint16

	for {
		pkt := <-stack.icmp4RecvCh
		peer := pkt.Peer
		payload := pkt.Payload

		if len(peer) != 4 {
			if stack.Verbose {
				log.Printf("icmp4 recv: drop(proto) %v", peer)
			}
			continue
		}

		if len(payload) < 4 {
			if stack.Verbose {
				log.Printf("icmp4 recv: drop(len) %v", peer)
			}
			continue
		}

		typ := payload[0]
		switch typ {
		case ICMP4EchoReq:
			stack.icmp4EchoRecv(pkt)
		case ICMP4EchoRep:
			stack.icmp4EchoRepRecv(pkt)
		default:
			if stack.Verbose {
				log.Printf("icmp4 recv: drop(typ) %v", peer)
			}
		}
	}
}

func (stack *IPStack) icmp4EchoRecv(pkt *IPPkt) {
	if stack.Verbose {
		log.Printf("icmp4 echo req recv: %v", pkt.Peer)
	}

	// Notice: modify pkt.Payload
	payload := pkt.Payload
	payload[0] = ICMP4EchoRep
	binary.BigEndian.PutUint16(payload[2:4], 0)
	cksum := util.CheckSum(payload)
	binary.BigEndian.PutUint16(payload[2:4], cksum)

	pkt.IfIdx = -1
	pkt.PeerMAC = nil

	if stack.Verbose {
		log.Printf("icmp4 echo rep send: %v", pkt.Peer)
	}

	stack.Send(pkt)
}

func (stack *IPStack) icmp4EchoRepRecv(pkt *IPPkt) {
	if stack.Verbose {
		log.Printf("icmp4 echo rep recv: %v", pkt.Peer)
	}
}
