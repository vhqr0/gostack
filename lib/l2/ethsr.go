package l2

import (
	"encoding/binary"
	"log"
	"net"
)

func (stack *EthStack) ethSender() {
	for {
		pkt := <-stack.sendCh
		iface := stack.Host.Ifaces[pkt.IfIdx]
		peer := pkt.Peer
		proto := pkt.Proto
		payload := pkt.Payload

		if stack.Verbose {
			log.Printf("eth send: %v", peer)
		}

		payloadLen := len(payload)
		padLen := 48 - payloadLen

		if padLen < 0 {
			padLen = 0
		}

		buf := make([]byte, 14+payloadLen+padLen)
		copy(buf[:6], peer)                           // dst
		copy(buf[6:12], iface.MAC)                    // src
		binary.BigEndian.PutUint16(buf[12:14], proto) // proto
		copy(buf[14:], payload)                       // payload

		if _, err := iface.Write(buf); err != nil {
			log.Panic(err)
		}
	}
}

func (stack *EthStack) ethReceiver(ifidx int) {
	// dst   byte[6]
	// src   byte[6]
	// proto uint16

	iface := stack.Host.Ifaces[ifidx]

	for {
		buf := make([]byte, 4096)
		if n, err := iface.Read(buf); err != nil {
			log.Panic(err)
		} else {
			buf = buf[:n]
		}

		if len(buf) < 14 {
			if stack.Verbose {
				log.Print("eth recv: drop(len)")
			}
			continue
		}

		dst := net.HardwareAddr(buf[:6])
		src := net.HardwareAddr(buf[6:12])
		proto := binary.BigEndian.Uint16(buf[12:14])
		payload := buf[14:]

		if _, ok := stack.MACSet[MACSetKey(ifidx, dst)]; !ok {
			if stack.Verbose {
				log.Printf("eth recv: drop(dst) %v => %v", src, dst)
			}
			continue
		}

		pkt := &EthPkt{
			IfIdx:   ifidx,
			Peer:    src,
			Proto:   proto,
			Payload: payload,
		}

		if ch, ok := stack.recvChMap[proto]; ok {
			if stack.Verbose {
				log.Printf("eth recv: %v => %v", src, dst)
			}
			ch <- pkt
		} else {
			if stack.Verbose {
				log.Printf("eth recv: drop(proto) %v => %v", src, dst)
			}
		}
	}
}
