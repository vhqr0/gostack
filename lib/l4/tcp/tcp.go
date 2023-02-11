package tcp

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/vhqr0/gostack/lib/l3"
	"github.com/vhqr0/gostack/lib/sock"
)

const (
	TCPFin = 0x01
	TCPSyn = 0x02
	TCPRst = 0x04
	TCPPsh = 0x08
	TCPAck = 0x10
	TCPUrg = 0x20
	TCPEcn = 0x40
	TCPWin = 0x80

	TCPOptNoop = 1
	TCPOptMss  = 2
	TCPOptSack = 5
	TCPOptTs   = 8

	TCPOptlenMss  = 4
	TCPOptlenSack = 2
)

type TCPPkt struct {
	Local *sock.Addr
	Peer  *sock.Addr

	SeqNum uint32
	AckNum uint32
	Fin    bool
	Syn    bool
	Rst    bool
	Ack    bool

	Payload []byte
}

func (pkt *TCPPkt) log(prefix string, from, to *sock.Addr) {
	info := fmt.Sprintf("%s %v => %v seq=%d", prefix, from, to, pkt.SeqNum)
	if pkt.Ack {
		info += fmt.Sprintf(" ack=%d", pkt.AckNum)
	}
	if pkt.Fin {
		info += " fin"
	}
	if pkt.Syn {
		info += " syn"
	}
	if pkt.Rst {
		info += " rst"
	}
	log.Print(info)
}

func (stack *TCPStack) tcpReceiver() {
	// sport      uint16
	// dport      uint16
	// seqnum     uint32
	// acknum     uint32
	// hlen|flags uint16
	// wsize      uint16
	// cksum      uint16
	// urgptr     uint16

	for {
		pkt := <-stack.recvCh
		payload := pkt.Payload

		if len(payload) < 20 {
			if stack.Verbose {
				log.Printf("tcp recv: drop(len) %v => %v", pkt.Peer, pkt.Local)
			}
			continue
		}

		local := &sock.Addr{IP: pkt.Local, Port: binary.BigEndian.Uint16(payload[2:4])}
		peer := &sock.Addr{IP: pkt.Peer, Port: binary.BigEndian.Uint16(payload[:2])}
		seqnum := binary.BigEndian.Uint32(payload[4:8])
		acknum := binary.BigEndian.Uint32(payload[8:12])
		flags := binary.BigEndian.Uint16(payload[12:14])
		hlen := (flags & 0xf000) >> 10
		fin := (flags & TCPFin) != 0
		syn := (flags & TCPSyn) != 0
		rst := (flags & TCPRst) != 0
		ack := (flags & TCPAck) != 0

		if len(payload) < int(hlen) {
			if stack.Verbose {
				log.Printf("tcp recv: drop(hdr) %v => %v", pkt.Peer, pkt.Local)
			}
			continue
		}

		tcppkt := &TCPPkt{
			Local:   local,
			Peer:    peer,
			SeqNum:  seqnum,
			AckNum:  acknum,
			Fin:     fin,
			Syn:     syn,
			Rst:     rst,
			Ack:     ack,
			Payload: payload[hlen:],
		}

		var sock *TCPSock

		if syn && !ack { // connect
			sock = stack.SockTable.LookupListen(local)
		} else {
			sock = stack.SockTable.LookupConn(local, peer)
		}

		if sock == nil {
			if stack.Verbose {
				tcppkt.log("tcp recv: drop(dst)", peer, local)
			}
			if !rst && (peer.IP.IsLoopback() || peer.IP.IsGlobalUnicast()) {
				tcppkt.SeqNum = 0
				tcppkt.AckNum = 0
				tcppkt.Fin = false
				tcppkt.Syn = false
				tcppkt.Rst = true
				tcppkt.Ack = false
				tcppkt.Payload = nil
				go stack.tcpSend(tcppkt)
			}
		} else {
			sock.tcpRecv(tcppkt)
		}
	}
}

func (stack *TCPStack) tcpSend(pkt *TCPPkt) error {
	local := pkt.Local
	peer := pkt.Peer
	seqnum := pkt.SeqNum
	acknum := pkt.AckNum
	fin := pkt.Fin
	syn := pkt.Syn
	rst := pkt.Rst
	ack := pkt.Ack
	payload := pkt.Payload

	flags := uint16(0x5000)
	if fin {
		flags |= TCPFin
	}
	if syn {
		flags |= TCPSyn
	}
	if rst {
		flags |= TCPRst
	}
	if ack {
		flags |= TCPAck
	}

	if stack.Verbose {
		pkt.log("tcp send:", local, peer)
	}

	buf := make([]byte, 20+len(payload))
	binary.BigEndian.PutUint16(buf[:2], local.Port) // sport
	binary.BigEndian.PutUint16(buf[2:4], peer.Port) // dport
	binary.BigEndian.PutUint32(buf[4:8], seqnum)    // seqnum
	binary.BigEndian.PutUint32(buf[8:12], acknum)   // acknum
	binary.BigEndian.PutUint16(buf[12:14], flags)   // hlen|flags
	binary.BigEndian.PutUint16(buf[14:16], 4096)    // wsize
	copy(buf[20:], payload)
	cksum := l3.InetCheckSum(buf, local.IP, peer.IP, l3.IPTCP)
	binary.BigEndian.PutUint16(buf[16:18], cksum)

	ippkt := &l3.IPPkt{
		IfIdx:   -1,
		Local:   local.IP,
		Peer:    peer.IP,
		Proto:   l3.IPTCP,
		Payload: buf,
	}

	return stack.IPStack.IPSend(ippkt) // Notice: block
}

func (sock *TCPSock) tcpRecv(pkt *TCPPkt) {}
