package conf

import (
	"github.com/vhqr0/gostack/lib/host"
)

type HostConf struct {
	Verbose bool
	Forward bool
	Ifaces  []*IfaceConf
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

func ConfFromHost(vhost *host.Host) *HostConf {
	conf := &HostConf{
		Verbose: vhost.Verbose,
		Forward: vhost.Forward,
	}
	for _, iface := range vhost.Ifaces {
		conf.Ifaces = append(conf.Ifaces, ConfFromIface(iface))
	}
	return conf
}
