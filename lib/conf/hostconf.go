package conf

import (
	"log"
	"net"

	"github.com/vhqr0/gostack/lib/host"
)

type HostConf struct {
	Verbose bool
	Forward bool
	Ifaces  []*IfaceConf
}

type IfaceConf struct {
	Name  string
	MTU   int
	MAC   string
	CIDR4 string
	CIDR6 string
	Typ   string
	Args  map[string]string
}

func (conf *HostConf) NewHost() *host.Host {
	vhost := &host.Host{
		Verbose: conf.Verbose,
		Forward: conf.Forward,
	}
	for _, ifaceConf := range conf.Ifaces {
		vhost.AddIface(ifaceConf.NewIface())
	}
	return vhost
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
