package l3

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/vhqr0/gostack/lib/l2"
	"github.com/vhqr0/gostack/lib/util"
)

const (
	NdpNS     uint8 = 135
	NdpNA     uint8 = 136
	NdpOptSrc uint8 = 1
	NdpOptTgt uint8 = 2
)

func (stack *IPStack) ndpRecv(pkt *IPPkt) {
	// typ    uint8    = NS | NA
	// code   uint8
	// cksum  uint16
	// flags  uint8    = 0xe0(RSO)
	// rsv    byte[3]
	// tgt    byte[16]
	// opttyp uint8    = Src | Tgt
	// optlen uint8    = 1
	// lladdr byte[6]

	payload := pkt.Payload

	if len(payload) < 32 {
		if stack.Verbose {
			log.Print("ndp recv: drop(len)")
		}
		return
	}

	typ := payload[0]
	query := net.IP(payload[8:24])

	if !query.IsGlobalUnicast() {
		if stack.Verbose {
			log.Printf("ndp recv: drop(query) %v", query)
		}
		return
	}

	switch typ {
	case NdpNS:
		stack.ndpNSRecv(pkt)
	case NdpNA:
		stack.ndpNARecv(pkt)
	default:
		if stack.Verbose {
			log.Print("ndp recv: drop(typ)")
		}
	}
}

func (stack *IPStack) ndpNSRecv(pkt *IPPkt) {
	ifidx := pkt.IfIdx
	peer := pkt.Peer
	payload := pkt.Payload

	iface := stack.EthStack.Host.Ifaces[ifidx]

	query := net.IP(payload[8:24])
	opttyp := payload[24]
	optlen := payload[25]
	peerMAC := net.HardwareAddr(payload[26:32])

	if opttyp != NdpOptSrc || optlen != 1 {
		if stack.Verbose {
			log.Printf("ndp ns recv: drop(opt) %v query %v", peerMAC, query)
		}
		return
	}

	if _, ok := stack.IPSet[string(query)]; !ok {
		if stack.Verbose {
			log.Printf("ndp ns recv: drop(dst) %v query %v", peerMAC, query)
		}
		return
	}

	if stack.Verbose {
		log.Printf("ndp ns recv: %v query %v", peerMAC, query)
	}

	buf := make([]byte, 72)
	buf[0] = 0x60               // ver
	buf[5] = 32                 // plen
	buf[6] = IPICMP6            // nh
	buf[7] = 0xff               // hlim
	copy(buf[8:24], query)      // src
	copy(buf[24:40], peer)      // dst
	buf[40] = NdpNA             // typ
	buf[44] = 0xe0              // flags
	copy(buf[48:64], query)     // tgt
	buf[64] = NdpOptTgt         // opttyp
	buf[65] = 1                 // optlen
	copy(buf[66:72], iface.MAC) // tgtmac
	cksum := util.CheckSum6(buf[40:72], query, peer, IPICMP6)
	binary.BigEndian.PutUint16(buf[42:44], cksum)

	ethPkt := &l2.EthPkt{
		IfIdx:   ifidx,
		Peer:    peerMAC,
		Proto:   l2.EthIP6,
		Payload: buf,
	}

	if stack.Verbose {
		log.Printf("ndp na send: %v answer %v", peerMAC, query)
	}

	stack.EthStack.Send(ethPkt)
}

func (stack *IPStack) ndpNARecv(pkt *IPPkt) {
	ifidx := pkt.IfIdx
	payload := pkt.Payload

	query := net.IP(payload[8:24])
	opttyp := payload[24]
	optlen := payload[25]
	peerMAC := net.HardwareAddr(payload[26:32])

	if opttyp != NdpOptTgt || optlen != 1 {
		if stack.Verbose {
			log.Printf("ndp na recv: drop(opt) %v answer %v", peerMAC, query)
		}
		return
	}

	if stack.Verbose {
		log.Printf("ndp na recv: %v answer %v", peerMAC, query)
	}

	stack.NeighTable.Update(Neigh6Key(ifidx, query), peerMAC)
}

func (stack *IPStack) ndpNSSend(ifidx int, peer net.IP) {
	if stack.Verbose {
		log.Printf("ndp ns send: query %v", peer)
	}

	iface := stack.EthStack.Host.Ifaces[ifidx]

	sol := util.SolIP(peer)

	buf := make([]byte, 40+32)
	buf[0] = 0x60               // ver
	buf[5] = 32                 // plen
	buf[6] = IPICMP6            // nh
	buf[7] = 0xff               // hlim
	copy(buf[8:24], iface.IP6)  // src
	copy(buf[24:40], sol)       // dst
	buf[40] = NdpNS             // typ
	copy(buf[48:64], peer)      // tgt
	buf[64] = NdpOptSrc         // opttyp
	buf[65] = 1                 // optlen
	copy(buf[66:72], iface.MAC) // srcmac
	cksum := util.CheckSum6(buf[40:72], iface.IP6, sol, IPICMP6)
	binary.BigEndian.PutUint16(buf[42:44], cksum)

	pkt := &l2.EthPkt{
		IfIdx:   ifidx,
		Peer:    util.SolMAC(peer),
		Proto:   l2.EthIP6,
		Payload: buf,
	}

	stack.EthStack.Send(pkt)
}
