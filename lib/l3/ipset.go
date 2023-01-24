package l3

import (
	"github.com/vhqr0/gostack/lib/host"
)

func NewIPSet(host *host.Host) map[string]struct{} {
	set := make(map[string]struct{})
	for _, iface := range host.Ifaces {
		set[string(iface.IP4)] = struct{}{}
		set[string(iface.IP6)] = struct{}{}
		set[string(SolIP(iface.IP6))] = struct{}{} // ip6 lookup
	}
	return set
}
