package tcp

import (
	"encoding/binary"
	"log"

	"github.com/vhqr0/gostack/lib/sock"
)

func TCPListenKey(addr *sock.Addr) (s string) {
	switch len(addr.IP) {
	case 4:
		s = TCP4ListenKey(addr)
	case 16:
		s = TCP6ListenKey(addr)
	default:
		log.Panic(&sock.InvalidSockAddr{Addr: addr})
	}
	return s
}

func TCPConnKey(local, peer *sock.Addr) (s string) {
	switch len(local.IP) {
	case 4:
		s = TCP4ConnKey(local, peer)
	case 16:
		s = TCP6ConnKey(local, peer)
	default:
		log.Panic(&sock.InvalidSockAddr{Addr: local})
	}
	return s
}

func TCP4ListenKey(addr *sock.Addr) string {
	buf := make([]byte, 6)
	copy(buf[:4], addr.IP)
	binary.BigEndian.PutUint16(buf[4:], addr.Port)
	return string(buf)
}

func TCP6ListenKey(addr *sock.Addr) string {
	buf := make([]byte, 18)
	copy(buf[:16], addr.IP)
	binary.BigEndian.PutUint16(buf[16:], addr.Port)
	return string(buf)
}

func TCP4ConnKey(local, peer *sock.Addr) string {
	buf := make([]byte, 12)
	copy(buf[:4], local.IP)
	binary.BigEndian.PutUint16(buf[4:6], local.Port)
	copy(buf[6:10], peer.IP)
	binary.BigEndian.PutUint16(buf[10:], peer.Port)
	return string(buf)
}

func TCP6ConnKey(local, peer *sock.Addr) string {
	buf := make([]byte, 36)
	copy(buf[:16], local.IP)
	binary.BigEndian.PutUint16(buf[16:18], local.Port)
	copy(buf[18:34], peer.IP)
	binary.BigEndian.PutUint16(buf[34:], peer.Port)
	return string(buf)
}
