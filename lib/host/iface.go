package host

import (
	"net"

	"github.com/vhqr0/gostack/lib/util"
)

type Iface struct {
	Name string
	MTU  int
	MAC  net.HardwareAddr
	IP4  net.IP
	IP6  net.IP
	Net4 *net.IPNet
	Net6 *net.IPNet

	adaptor IfaceAdaptor
}

func (iface *Iface) Read(p []byte) (int, error) {
	return iface.adaptor.Read(p)
}

func (iface *Iface) Write(p []byte) (int, error) {
	return iface.adaptor.Write(p)
}

func (iface *Iface) Close() error {
	return iface.adaptor.Close()
}

func NewIface(name, cidr4, cidr6, typ string, args map[string]string) (*Iface, error) {
	adaptor, err := NewIfaceAdaptor(typ, args)
	if err != nil {
		return nil, err
	}

	ip4, net4, err := util.ParseCIDR4(cidr4)
	if err != nil {
		return nil, err
	}

	ip6, net6, err := util.ParseCIDR6(cidr6)
	if err != nil {
		return nil, err
	}

	mac := net.HardwareAddr(util.RandBytes(6))

	iface := &Iface{
		Name: name,
		MTU:  1500,
		MAC:  mac,
		IP4:  ip4,
		IP6:  ip6,
		Net4: net4,
		Net6: net6,

		adaptor: adaptor,
	}

	return iface, nil
}
