package l3

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/vhqr0/gostack/lib/l2"
)

const (
	IP6Frag uint8 = 44
)

func (stack *IPStack) ip6Receiver() {
	// ver|prio uint8
	// fl       byte[3]
	// plen     uint16
	// nh       uint8
	// hlimit   uint8
	// src      byte[16]
	// dst      byte[16]

	for {
		pkt := <-stack.ip6RecvCh
		ifidx := pkt.IfIdx
		peerMAC := pkt.Peer
		payload := pkt.Payload

		if len(payload) < 40 {
			if stack.Verbose {
				log.Print("ip6 recv: drop(len)")
			}
			continue
		}

		ver := (payload[0] & 0xf0) >> 4
		plen := binary.BigEndian.Uint16(payload[4:6])
		nh := payload[6]
		src := net.IP(payload[8:24])
		dst := net.IP(payload[24:40])

		if ver != 6 || len(payload) < int(40+plen) {
			if stack.Verbose {
				log.Print("ip6 recv: drop(hdr)")
			}
			continue
		}

		// restrict:
		// src: global unicast
		// dst: global unicast or multicast (in ipset)

		if !src.IsGlobalUnicast() {
			if stack.Verbose {
				log.Printf("ip6 recv: drop(src) %v => %v", src, dst)
			}
		}

		if !dst.IsGlobalUnicast() && !dst.IsMulticast() {
			if stack.Verbose {
				log.Printf("ip6 recv: drop(dst) %v => %v", src, dst)
			}
		}

		if _, ok := stack.IPSet[string(dst)]; !ok {
			if stack.Forward && !dst.IsMulticast() {
				if stack.Verbose {
					log.Printf("ip6 recv: forward %v => %v", src, dst)
				}
				go stack.ipForward(pkt, dst) // Notice: go, or maybe block if not
			} else if stack.Verbose {
				log.Printf("ip6 recv: drop(dst) %v => %v", src, dst)
			}
		}

		if nh == IP6Frag { // TODO: handle frag
			if stack.Verbose {
				log.Printf("ip6 recv: drop(frag) %v => %v", src, dst)
			}
			continue
		}

		ippkt := &IPPkt{
			IfIdx:   ifidx,
			PeerMAC: peerMAC,
			Local:   dst,
			Peer:    src,
			Proto:   nh,
			Payload: payload[40 : 40+plen],
		}

		stack.ipRecv(ippkt)
	}
}

func (stack *IPStack) ip6Send(pkt *IPPkt) error {
	ifidx := pkt.IfIdx
	peerMAC := pkt.PeerMAC
	local := pkt.Local
	peer := pkt.Peer
	proto := pkt.Proto
	payload := pkt.Payload

	iface := stack.EthStack.Host.Ifaces[ifidx]

	plen := len(payload)

	if 40+plen > iface.MTU {
		// TODO: ipid, frag
		if stack.Verbose {
			log.Printf("ip6 send: drop(ptb) %v => %v", local, peer)
		}
		return &PktTooBig{Len: 40 + plen}
	}

	buf := make([]byte, 40+plen)
	buf[0] = 0x60                                      // magic
	binary.BigEndian.PutUint16(buf[4:6], uint16(plen)) // plen
	buf[6] = proto                                     // nh
	buf[7] = 0xff                                      // hlim
	copy(buf[8:24], local)                             // src
	copy(buf[24:40], peer)                             // dst
	copy(buf[40:], payload)                            // payload

	ethPkt := &l2.EthPkt{
		IfIdx:   ifidx,
		Peer:    peerMAC,
		Proto:   l2.EthIP6,
		Payload: buf,
	}

	stack.EthStack.Send(ethPkt)

	return nil
}
