package udp

import (
	"encoding/binary"
	"log"

	"github.com/vhqr0/gostack/lib/sock"
)

func UDPKey(addr *sock.Addr) (s string) {
	switch len(addr.IP) {
	case 4:
		s = UDP4Key(addr)
	case 16:
		s = UDP6Key(addr)
	default:
		log.Panic(&sock.InvalidSockAddr{Addr: addr})
	}
	return s
}

func UDP4Key(addr *sock.Addr) string {
	buf := make([]byte, 6)
	copy(buf[:4], addr.IP)
	binary.BigEndian.PutUint16(buf[4:], addr.Port)
	return string(buf)
}

func UDP6Key(addr *sock.Addr) string {
	buf := make([]byte, 18)
	copy(buf[:16], addr.IP)
	binary.BigEndian.PutUint16(buf[16:], addr.Port)
	return string(buf)
}
