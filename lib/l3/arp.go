package l3

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"github.com/vhqr0/gostack/lib/l2"
	"github.com/vhqr0/gostack/lib/util"
)

const (
	ARPEth uint16 = 0x0001
	ARPIP4 uint16 = 0x0800
	ARPReq uint16 = 0x0001
	ARPRep uint16 = 0x0002
)

func (stack *IPStack) arpReceiver() {
	// hwtyp   uint16  = eth
	// protyp  uint16  = ip4
	// hwsize  uint8   = 6
	// prosize uint8   = 4
	// opcode  uint16  = req | rep
	// hwsrc   byte[6]
	// ipsrc   byte[4]
	// hwdst   byte[6]
	// ipdst   byte[4]

	for {
		pkt := <-stack.arpRecvCh
		payload := pkt.Payload

		if len(payload) < 28 {
			if stack.Verbose {
				log.Print("arp recv: drop(len)")
			}
			continue
		}

		hwtyp := binary.BigEndian.Uint16(payload[:2])
		protyp := binary.BigEndian.Uint16(payload[2:4])
		hwsize := payload[4]
		prosize := payload[5]
		opcode := binary.BigEndian.Uint16(payload[6:8])

		if hwtyp != ARPEth || protyp != ARPIP4 || hwsize != 6 || prosize != 4 {
			if stack.Verbose {
				log.Print("arp recv: drop(hdr)")
			}
			continue
		}

		switch opcode {
		case ARPReq:
			stack.arpReqRecv(pkt)
		case ARPRep:
			stack.arpRepRecv(pkt)
		default:
			if stack.Verbose {
				log.Printf("arp recv: drop(op)")
			}
		}
	}
}

func (stack *IPStack) arpReqRecv(pkt *l2.EthPkt) {
	payload := pkt.Payload

	iface := stack.EthStack.Host.Ifaces[pkt.IfIdx]

	// Notice: copy before modify pkt.Payload
	query := net.IP(make([]byte, 4))
	copy(query, payload[24:28])

	if _, ok := stack.IPSet[string(query)]; !ok {
		if stack.Verbose {
			log.Printf("arp req recv: drop(dst) %v query %v", net.HardwareAddr(payload[8:14]), query)
		}
		return
	}

	if stack.Verbose {
		log.Printf("arp req recv: %v query %v", net.HardwareAddr(payload[8:14]), query)
	}

	// Notice: modify pkt.Payload
	binary.BigEndian.PutUint16(payload[6:8], ARPRep) // toggle opcode
	copy(payload[18:28], payload[8:18])              // copy src to dst
	copy(payload[8:14], iface.MAC)                   // hwsrc: iface.src
	copy(payload[14:18], query)                      // ipsrc: query

	if stack.Verbose {
		log.Printf("arp rep send: %v answer %v", pkt.Peer, query)
	}

	stack.EthStack.Send(pkt)
}

func (stack *IPStack) arpRepRecv(pkt *l2.EthPkt) {
	ifidx := pkt.IfIdx
	payload := pkt.Payload

	iface := stack.EthStack.Host.Ifaces[ifidx]

	src := net.HardwareAddr(payload[8:14])
	dst := net.HardwareAddr(payload[18:24])
	query := net.IP(payload[14:18])

	if !bytes.Equal(dst, iface.MAC) {
		if stack.Verbose {
			log.Printf("arp rep recv: drop(dst) %v answer %v", src, dst)
		}
		return
	}

	if stack.Verbose {
		log.Printf("arp rep recv: %v answer %v", src, query)
	}

	stack.NeighTable.Update(Neigh4Key(ifidx, query), src)
}

func (stack *IPStack) arpReqSend(ifidx int, peer net.IP) {
	if stack.Verbose {
		log.Printf("arp req send: query %v", peer)
	}

	iface := stack.EthStack.Host.Ifaces[ifidx]

	payload := make([]byte, 28)
	binary.BigEndian.PutUint16(payload[:2], ARPEth)
	binary.BigEndian.PutUint16(payload[2:4], ARPIP4)
	payload[4] = 6
	payload[5] = 4
	binary.BigEndian.PutUint16(payload[6:8], ARPReq)
	copy(payload[8:14], iface.MAC)
	copy(payload[14:18], iface.IP4)
	copy(payload[24:28], peer)

	pkt := &l2.EthPkt{
		IfIdx:   ifidx,
		Peer:    util.BrdMAC,
		Proto:   l2.EthARP,
		Payload: payload,
	}

	stack.EthStack.Send(pkt)
}
