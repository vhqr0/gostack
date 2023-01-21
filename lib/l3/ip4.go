package l3

import (
	"encoding/binary"
	"log"
	"math/rand"
	"net"

	"github.com/vhqr0/gostack/lib/l2"
	"github.com/vhqr0/gostack/lib/util"
)

func (stack *IPStack) ip4Receiver() {
	// ver|hlen     uint8  = 0x45
	// tos          uint8
	// tlen         uint16
	// id           uint16
	// flags|offset uint16
	// ttl          uint8
	// proto        uint8
	// cksum        uint16
	// src          byte[4]
	// dst          byte[4]

	for {
		pkt := <-stack.ip4RecvCh
		ifidx := pkt.IfIdx
		peerMAC := pkt.Peer
		payload := pkt.Payload

		if len(payload) < 20 {
			if stack.Verbose {
				log.Print("ip4 recv: drop(len)")
			}
			continue
		}

		magic := payload[0]
		tlen := binary.BigEndian.Uint16(payload[2:4])
		fragFlags := (payload[6] & 0xe0) >> 5
		fragOffset := binary.BigEndian.Uint16(payload[6:8]) & 0x1fff
		proto := payload[9]
		src := net.IP(payload[12:16])
		dst := net.IP(payload[16:20])

		if magic != 0x45 || len(payload) < int(tlen) {
			if stack.Verbose {
				log.Print("ip4 recv: drop(hdr)")
			}
			continue
		}

		// restrict:
		// src: global unicast
		// dst: global unicast

		if !src.IsGlobalUnicast() {
			if stack.Verbose {
				log.Printf("ip4 recv: drop(src) %v => %v", src, dst)
			}
			continue
		}

		if !dst.IsGlobalUnicast() {
			if stack.Verbose {
				log.Printf("ip4 recv: drop(dst) %v => %v", src, dst)
			}
			continue
		}

		if _, ok := stack.IPSet[string(dst)]; !ok {
			if stack.Forward {
				if stack.Verbose {
					log.Printf("ip4 recv: forward %v => %v", src, dst)
				}
				go stack.ipForward(pkt, dst) // Notice: go, or maybe block if not
			} else if stack.Verbose {
				log.Printf("ip4 recv: drop(dst) %v => %v", src, dst)
			}
			continue
		}

		if fragFlags&1 != 0 || fragOffset != 0 { // TODO: handle frag
			if stack.Verbose {
				log.Printf("ip4 recv: drop(frag) %v => %v", src, dst)
			}
			continue
		}

		ippkt := &IPPkt{
			IfIdx:   ifidx,
			PeerMAC: peerMAC,
			Local:   dst,
			Peer:    src,
			Proto:   proto,
			Payload: payload[20:tlen],
		}

		stack.ipRecv(ippkt)
	}
}

func (stack *IPStack) ip4Send(pkt *IPPkt) error {
	ifidx := pkt.IfIdx
	peerMAC := pkt.PeerMAC
	local := pkt.Local
	peer := pkt.Peer
	proto := pkt.Proto
	payload := pkt.Payload

	iface := stack.EthStack.Host.Ifaces[ifidx]

	tlen := 20 + len(payload)

	if tlen > iface.MTU {
		// TODO: frag
		if stack.Verbose {
			log.Printf("ip4 send: drop(ptb) %v => %v", local, peer)
		}
		return &PktTooBig{Len: tlen}
	}

	buf := make([]byte, tlen)
	buf[0] = 0x45                                               // magic
	binary.BigEndian.PutUint16(buf[2:4], uint16(tlen))          // tlen
	binary.BigEndian.PutUint16(buf[4:6], uint16(rand.Uint32())) // id
	buf[8] = 0xff                                               // ttl
	buf[9] = proto                                              // proto
	copy(buf[12:16], local)                                     // src
	copy(buf[16:20], peer)                                      // dst
	copy(buf[20:], payload)                                     // payload
	cksum := util.CheckSum(buf)
	binary.BigEndian.PutUint16(buf[10:12], cksum)

	ethPkt := &l2.EthPkt{
		IfIdx:   ifidx,
		Peer:    peerMAC,
		Proto:   l2.EthIP4,
		Payload: buf,
	}

	stack.EthStack.Send(ethPkt)

	return nil
}
