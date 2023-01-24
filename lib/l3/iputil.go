package l3

import (
	"log"
	"net"

	"github.com/vhqr0/gostack/lib/util"
)

var (
	SolIPPrefix = []byte("\xff\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\xff")
)

func InetCheckSum(msg []byte, src, dst net.IP, proto uint8) (cksum uint16) {
	switch len(src) {
	case 4:
		cksum = util.Inet4CheckSum(msg, src, dst, proto)
	case 16:
		cksum = util.Inet6CheckSum(msg, src, dst, proto)
	default:
		log.Panic(&InvalidIPLen{Len: len(src)})
	}
	return cksum
}

func AllIPNet(ver int) (*net.IPNet, error) {
	ipnet := &net.IPNet{}
	switch ver {
	case 4:
		ipnet.IP = net.IPv4zero.To4()
		ipnet.Mask = net.CIDRMask(0, 32)
	case 6:
		ipnet.IP = net.IPv6zero
		ipnet.Mask = net.CIDRMask(0, 128)
	default:
		return nil, &InvalidIPVer{Ver: ver}
	}
	return ipnet, nil
}

func ParseIP(ver int, ipStr string) (net.IP, error) {
	switch ver {
	case 4:
		return util.ParseIP4(ipStr)
	case 6:
		return util.ParseIP6(ipStr)
	}
	return nil, &InvalidIPVer{Ver: ver}
}

func ParseCIDR(ver int, cidr string) (net.IP, *net.IPNet, error) {
	switch ver {
	case 4:
		return util.ParseCIDR4(cidr)
	case 6:
		return util.ParseCIDR6(cidr)
	}
	return nil, nil, &InvalidIPVer{Ver: ver}
}

func SolIP(ip net.IP) net.IP {
	sip := make([]byte, 16)
	copy(sip[:13], SolIPPrefix)
	copy(sip[13:], ip[13:])
	return sip
}
