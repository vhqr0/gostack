package l3

import (
	"encoding/binary"
	"log"
	"net"
)

func NeighKey(ifidx int, ip net.IP) (key string) {
	switch len(ip) {
	case 4:
		key = Neigh4Key(ifidx, ip)
	case 16:
		key = Neigh6Key(ifidx, ip)
	default:
		log.Panic(&InvalidIPLen{Len: len(ip)})
	}
	return
}

func Neigh4Key(ifidx int, ip net.IP) string {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[:2], uint16(ifidx))
	copy(buf[2:6], ip)
	return string(buf)
}

func Neigh6Key(ifidx int, ip net.IP) string {
	buf := make([]byte, 18)
	binary.BigEndian.PutUint16(buf[:2], uint16(ifidx))
	copy(buf[2:18], ip)
	return string(buf)
}
