package sock

import (
	"net"
	"strconv"
)

type Addr struct {
	IP   net.IP
	Port uint16
}

func (addr *Addr) Equal(o *Addr) bool {
	return addr.Port == o.Port && addr.IP.Equal(o.IP)
}

func (addr *Addr) String() string {
	if addr == nil {
		return "<nil>"
	}
	portStr := ":" + strconv.FormatUint(uint64(addr.Port), 10)
	switch len(addr.IP) {
	case 4:
		return addr.IP.String() + portStr
	case 16:
		return "[" + addr.IP.String() + "]" + portStr
	}
	return "<nil>" + portStr
}

func ResolveAddr(addrStr string) (addr *Addr, family uint32, err error) {
	var tcpAddr *net.TCPAddr
	if tcpAddr, err = net.ResolveTCPAddr("tcp", addrStr); err != nil {
		return
	}
	addr = &Addr{Port: uint16(tcpAddr.Port)}
	if ip4 := tcpAddr.IP.To4(); ip4 != nil {
		family = AFIP4
		addr.IP = ip4
	} else {
		family = AFIP6
		addr.IP = tcpAddr.IP
	}
	return
}
