package conf

import (
	"log"
	"net"

	"github.com/vhqr0/gostack/lib/host"
	"github.com/vhqr0/gostack/lib/util"
)

type IfaceConf struct {
	Name  string
	MTU   int
	MAC   string
	CIDR4 string
	CIDR6 string
	Typ   string
	Args  map[string]string
}

func (conf *IfaceConf) NewIface() *host.Iface {
	iface, err := host.NewIface(conf.Name, conf.CIDR4, conf.CIDR6, conf.Typ, conf.Args)
	if err != nil {
		log.Panic(err)
	}
	if conf.MTU != 0 {
		iface.MTU = conf.MTU
	}
	if conf.MAC != "" {
		if mac, err := net.ParseMAC(conf.MAC); err != nil {
			log.Panic(err)
		} else {
			iface.MAC = mac
		}
	}
	return iface
}

func IfaceConfFrom(iface *host.Iface) *IfaceConf {
	conf := &IfaceConf{
		Name: iface.Name,
		MTU:  iface.MTU,
		MAC:  iface.MAC.String(),
		CIDR4: util.CIDRString(iface.IP4, iface.Net4),
		CIDR6: util.CIDRString(iface.IP6, iface.Net6),
		Typ:  iface.Typ,
		Args: iface.Args,
	}
	return conf
}
