package util

import (
	"net"
)

var (
	BrdMAC       = net.HardwareAddr("\xff\xff\xff\xff\xff\xff")
	SolIPPrefix  = []byte("\xff\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\xff")
	SolMACPrefix = []byte("\x33\x33\xff")
)

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
		return ParseIP4(ipStr)
	case 6:
		return ParseIP6(ipStr)
	}
	return nil, &InvalidIPVer{Ver: ver}
}

func ParseCIDR(ver int, cidr string) (net.IP, *net.IPNet, error) {
	switch ver {
	case 4:
		return ParseCIDR4(cidr)
	case 6:
		return ParseCIDR6(cidr)
	}
	return nil, nil, &InvalidIPVer{Ver: ver}
}

func ParseIP4(ipStr string) (net.IP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, &InvalidIPStr{Str: ipStr}
	}
	if ip = ip.To4(); ip == nil {
		return nil, &InvalidIPStr{Str: ipStr}
	}
	if !ip.IsGlobalUnicast() {
		return nil, &InvalidIPStr{Str: ipStr}
	}
	return ip, nil
}

func ParseIP6(ipStr string) (net.IP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, &InvalidIPStr{Str: ipStr}
	}
	if ip = ip.To16(); ip == nil {
		return nil, &InvalidIPStr{Str: ipStr}
	}
	if !ip.IsGlobalUnicast() {
		return nil, &InvalidIPStr{Str: ipStr}
	}
	return ip, nil
}

func ParseCIDR4(cidr string) (net.IP, *net.IPNet, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, err
	}
	if _, bits := ipnet.Mask.Size(); bits != 32 {
		return nil, nil, &InvalidIPStr{Str: cidr}
	}
	if ip = ip.To4(); ip == nil {
		return nil, nil, &InvalidIPStr{Str: cidr}
	}
	if !ip.IsGlobalUnicast() {
		return nil, nil, &InvalidIPStr{Str: cidr}
	}
	return ip, ipnet, nil
}

func ParseCIDR6(cidr string) (net.IP, *net.IPNet, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, err
	}
	if _, bits := ipnet.Mask.Size(); bits != 128 {
		return nil, nil, &InvalidIPStr{Str: cidr}
	}
	if ip = ip.To16(); ip == nil {
		return nil, nil, &InvalidIPStr{Str: cidr}
	}
	if !ip.IsGlobalUnicast() {
		return nil, nil, &InvalidIPStr{Str: cidr}
	}
	return ip, ipnet, nil
}

func CIDRString(ip net.IP, ipnet *net.IPNet) string {
	cidr := &net.IPNet{
		IP: ip,
		Mask: ipnet.Mask,
	}
	return cidr.String()
}

func SolIP(ip net.IP) net.IP {
	sip := make([]byte, 16)
	copy(sip[:13], SolIPPrefix)
	copy(sip[13:], ip[13:])
	return sip
}

func SolMAC(ip net.IP) net.HardwareAddr {
	smac := make([]byte, 6)
	smac[0] = 0x33
	smac[1] = 0x33
	smac[2] = 0xff
	copy(smac[3:], ip[13:])
	return smac
}
