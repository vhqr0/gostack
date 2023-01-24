package util

import (
	"encoding/binary"
	"net"
)

func CheckSum(msg []byte) uint16 {
	sum := uint32(0)
	var i int
	for i = 0; i < len(msg)-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(msg[i : i+2]))
	}
	if i == len(msg)-1 { // odd
		sum += uint32(uint16(msg[i]) << 8)
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return ^uint16(sum)
}

func Inet4CheckSum(msg []byte, src, dst net.IP, proto uint8) uint16 {
	msgLen := len(msg)
	buf := make([]byte, 12+msgLen)
	copy(buf[:4], src)                                    // src
	copy(buf[4:8], dst)                                   // dst
	binary.BigEndian.PutUint16(buf[8:10], uint16(msgLen)) // plen
	binary.BigEndian.PutUint16(buf[10:12], uint16(proto)) // proto
	copy(buf[12:], msg)
	return CheckSum(buf)
}

func Inet6CheckSum(msg []byte, src, dst net.IP, proto uint8) uint16 {
	msgLen := len(msg)
	buf := make([]byte, 36+msgLen)
	copy(buf[:16], src)                                    // src
	copy(buf[16:32], dst)                                  // dst
	binary.BigEndian.PutUint16(buf[32:34], uint16(msgLen)) // plen
	binary.BigEndian.PutUint16(buf[34:36], uint16(proto))  // proto
	copy(buf[36:], msg)
	return CheckSum(buf)
}
