package l2

import (
	"encoding/binary"
	"net"
)

func MACSetKey(ifidx int, mac net.HardwareAddr) string {
	key := make([]byte, 8)
	binary.BigEndian.PutUint16(key[:2], uint16(ifidx))
	copy(key[2:], mac)
	return string(key)
}
