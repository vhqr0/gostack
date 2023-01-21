package l3

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/vhqr0/gostack/lib/l2"
	"github.com/vhqr0/gostack/lib/util"
)

func (stack *IPStack) ipSender() {
	for {
		pkt := <-stack.sendCh
		go stack.IPSend(pkt) // Notice: go, or maybe block if not
	}
}

func (stack *IPStack) IPSend(pkt *IPPkt) error { // Notice: block, return error, public
	ifidx := pkt.IfIdx
	peerMAC := pkt.PeerMAC
	local := pkt.Local
	peer := pkt.Peer

	// restrict:
	// peer: loopback or global unicast

	if stack.isInStack(peer) { // capture loopback
		if stack.Verbose {
			log.Printf("ip send: loopback %v => %v", local, peer)
		}
		pkt.IfIdx = -1
		pkt.PeerMAC = nil
		stack.ipRecv(pkt)
		return nil
	}

	if !peer.IsGlobalUnicast() {
		log.Printf("ip send: drop(dst) %v => %v", local, peer)
		return &DstUnreach{Dst: peer}
	}

	dstPeer := peer // for route

	if ifidx < 0 {
		entry := stack.Next(peer)
		if entry == nil {
			log.Printf("ip send: drop(dst) %v => %v", local, peer)
			return &DstUnreach{Dst: peer}
		}
		ifidx = entry.IfIdx
		if entry.Peer != nil {
			dstPeer = entry.Peer
		}
	}

	iface := stack.EthStack.Host.Ifaces[ifidx]

	if local == nil {
		switch len(peer) {
		case 4:
			local = iface.IP4
		case 16:
			local = iface.IP6
		default:
			log.Panic(&util.InvalidIPLen{Len: len(peer)})
		}
	}

	if peerMAC == nil {
		peerMAC = stack.Lookup(ifidx, dstPeer)
		if peerMAC == nil {
			log.Printf("ip send: drop(host) %v => %v", local, peer)
			return &HostUnreach{Host: dstPeer}
		}
	}

	pkt.IfIdx = ifidx
	pkt.PeerMAC = peerMAC
	pkt.Local = local
	pkt.Peer = peer

	if stack.Verbose {
		log.Printf("ip send: %v => %v", local, peer)
	}

	switch len(peer) {
	case 4:
		return stack.ip4Send(pkt)
	case 16:
		return stack.ip6Send(pkt)
	default:
		log.Panic(&util.InvalidIPLen{Len: len(peer)})
	}

	return nil
}

func (stack *IPStack) ipRecv(pkt *IPPkt) {
	if ch, ok := stack.recvChMap[pkt.Proto]; ok {
		if stack.Verbose {
			log.Printf("ip recv: %v => %v", pkt.Peer, pkt.Local)
		}
		ch <- pkt
	} else {
		if stack.Verbose {
			log.Printf("ip recv: drop(proto) %v => %v", pkt.Peer, pkt.Local)
		}
	}
}

func (stack *IPStack) ipForward(pkt *l2.EthPkt, dst net.IP) {
	entry := stack.Next(dst)
	if entry == nil {
		if stack.Verbose {
			log.Printf("ip forward: drop(dst) %v", dst)
		}
		return
	}

	ifidx := entry.IfIdx
	peer := entry.Peer
	if peer == nil {
		peer = dst
	}

	iface := stack.EthStack.Host.Ifaces[ifidx]
	if len(pkt.Payload) > iface.MTU {
		if stack.Verbose {
			log.Printf("ip forward: drop(ptb) %v via %v", dst, peer)
		}
		return
	}

	payload := pkt.Payload
	ttl := payload[9]
	if ttl == 1 {
		if stack.Verbose {
			log.Printf("ip forward: drop(ttl) %v via %v", dst, peer)
		}
		return
	}

	peerMAC := stack.Lookup(ifidx, peer)
	if peerMAC == nil {
		if stack.Verbose {
			log.Printf("ip forward: drop(host) %v via %v", dst, peer)
		}
		return
	}

	// Notice: modify pkt.Payload
	payload[9] = ttl - 1
	binary.BigEndian.PutUint16(payload[10:12], 0)
	cksum := util.CheckSum(payload)
	binary.BigEndian.PutUint16(payload[10:12], cksum)

	pkt.IfIdx = ifidx
	pkt.Peer = peerMAC

	if stack.Verbose {
		log.Printf("ip forward: %v via %v", dst, peer)
	}

	stack.EthStack.Send(pkt)
}

func (stack *IPStack) isInStack(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}
	if ip.IsMulticast() { // ipset contains multicast address
		return false
	}
	_, ok := stack.IPSet[string(ip)]
	return ok
}
