package l3

import (
	"encoding/binary"
	"log"

	"github.com/vhqr0/gostack/lib/util"
)

const (
	ICMP6EchoReq uint8 = 128
	ICMP6EchoRep uint8 = 129
)

func (stack *IPStack) icmp6Receiver() {
	// typ   uint8
	// code  uint8
	// cksum uint16

	for {
		pkt := <-stack.icmp6RecvCh
		peer := pkt.Peer
		payload := pkt.Payload

		if len(peer) != 16 {
			if stack.Verbose {
				log.Printf("icmp6 recv: drop(proto) %v", peer)
			}
			continue
		}

		if len(payload) < 4 {
			if stack.Verbose {
				log.Printf("icmp6 recv: drop(len) %v", peer)
			}
			continue
		}

		typ := payload[0]
		switch typ {
		case NdpNS, NdpNA:
			stack.ndpRecv(pkt)
		case ICMP6EchoReq:
			stack.icmp6EchoReqRecv(pkt)
		case ICMP6EchoRep:
			stack.icmp6EchoRepRecv(pkt)
		default:
			if stack.Verbose {
				log.Printf("icmp6 recv: drop(typ) %v", peer)
			}
		}
	}
}

func (stack *IPStack) icmp6EchoReqRecv(pkt *IPPkt) {
	if stack.Verbose {
		log.Printf("icmp6 echo req recv: %v", pkt.Peer)
	}

	// Notice: modify pkt.Payload
	payload := pkt.Payload
	payload[0] = ICMP6EchoRep
	binary.BigEndian.PutUint16(payload[2:4], 0)
	cksum := util.CheckSum6(payload, pkt.Local, pkt.Peer, IPICMP6)
	binary.BigEndian.PutUint16(payload[2:4], cksum)

	pkt.IfIdx = -1
	pkt.PeerMAC = nil

	if stack.Verbose {
		log.Printf("icmp6 echo rep send: %v", pkt.Peer)
	}

	stack.Send(pkt)
}

func (stack *IPStack) icmp6EchoRepRecv(pkt *IPPkt) {
	if stack.Verbose {
		log.Printf("icmp6 echo rep recv: %v", pkt.Peer)
	}
}
