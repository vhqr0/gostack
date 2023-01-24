package l2

import (
	"net"
)

var (
	BrdMAC      = net.HardwareAddr("\xff\xff\xff\xff\xff\xff")
	olMACPrefix = []byte("\x33\x33\xff")
)

func SolMAC(ip net.IP) net.HardwareAddr {
	smac := make([]byte, 6)
	smac[0] = 0x33
	smac[1] = 0x33
	smac[2] = 0xff
	copy(smac[3:], ip[13:])
	return smac
}
