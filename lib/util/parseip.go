package util

import (
	"net"
)

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
		IP:   ip,
		Mask: ipnet.Mask,
	}
	return cidr.String()
}
