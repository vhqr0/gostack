package udp

import (
	"encoding/binary"
	"log"

	"github.com/vhqr0/gostack/lib/l3"
	"github.com/vhqr0/gostack/lib/sock"
)

type UDPPkt struct {
	Local   *sock.Addr
	Peer    *sock.Addr
	Payload []byte
}

func (stack *UDPStack) udpReceiver() {
	// sport uint16
	// dport uint16
	// tlen  uint16
	// cksum uint16

	for {
		pkt := <-stack.recvCh
		payload := pkt.Payload

		if len(payload) < 8 {
			if stack.Verbose {
				log.Printf("udp recv: drop(len) %v => %v", pkt.Peer, pkt.Local)
			}
			continue
		}

		tlen := binary.BigEndian.Uint16(payload[4:6])
		local := &sock.Addr{IP: pkt.Local, Port: binary.BigEndian.Uint16(payload[2:4])}
		peer := &sock.Addr{IP: pkt.Peer, Port: binary.BigEndian.Uint16(payload[:2])}

		if len(payload) != int(tlen) {
			if stack.Verbose {
				log.Printf("udp recv: drop(hdr) %v => %v", pkt.Peer, pkt.Local)
			}
			continue
		}

		udppkt := &UDPPkt{
			Local:   local,
			Peer:    peer,
			Payload: payload[8:],
		}

		sock := stack.SockTable.Lookup(local)
		if sock == nil {
			if stack.Verbose {
				log.Printf("udp recv: drop(dst) %v => %v", peer, local)
			}
		} else {
			sock.udpRecv(udppkt)
		}
	}
}

func (stack *UDPStack) udpSend(pkt *UDPPkt) error { // Notice: block. return error
	local := pkt.Local
	peer := pkt.Peer
	payload := pkt.Payload

	if stack.Verbose {
		log.Printf("udp send: %v => %v", local, peer)
	}

	tlen := 8 + len(payload)

	buf := make([]byte, tlen)
	binary.BigEndian.PutUint16(buf[:2], local.Port)    // sport
	binary.BigEndian.PutUint16(buf[2:4], peer.Port)    // dport
	binary.BigEndian.PutUint16(buf[4:6], uint16(tlen)) // tlen
	copy(buf[8:], payload)                             // payload
	cksum := l3.InetCheckSum(buf, local.IP, peer.IP, l3.IPUDP)
	binary.BigEndian.PutUint16(buf[6:8], cksum)

	ippkt := &l3.IPPkt{
		IfIdx:   -1,
		Local:   local.IP,
		Peer:    peer.IP,
		Proto:   l3.IPUDP,
		Payload: buf,
	}

	return stack.IPStack.IPSend(ippkt) // Notice: block
}

func (sock *UDPSock) udpRecv(pkt *UDPPkt) {
	sock.Mutex.RLock()
	defer sock.Mutex.RUnlock()

	if sock.Status == UDPClosed {
		if sock.Stack.Verbose {
			log.Printf("udp recv: drop(closed) %v => %v", pkt.Peer, pkt.Local)
		}
		return
	}

	if sock.Peer != nil && !pkt.Peer.Equal(sock.Peer) {
		if sock.Stack.Verbose {
			log.Printf("udp recv: drop(src) %v => %v", pkt.Peer, pkt.Local)
		}
		return
	}

	select {
	case sock.recvCh <- pkt:
		if sock.Stack.Verbose {
			log.Printf("udp recv: %v => %v", pkt.Peer, pkt.Local)
		}
	default:
		if sock.Stack.Verbose {
			log.Printf("udp recv: drop(listen) %v => %v", pkt.Peer, pkt.Local)
		}
	}

}
