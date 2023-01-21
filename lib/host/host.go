package host

type Host struct {
	Verbose bool // inherit to * (l2, l3, l4)
	Forward bool // inherit to l3

	Ifaces []*Iface
}

func (host *Host) AddIface(iface *Iface) {
	host.Ifaces = append(host.Ifaces, iface)
}

func (host *Host) GetIfaceByName(ifname string) int {
	for ifidx, iface := range host.Ifaces {
		if iface.Name == ifname {
			return ifidx
		}
	}
	return -1
}
