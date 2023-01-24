package l2

import (
	"github.com/vhqr0/gostack/lib/host"
)

func NewMACSet(host *host.Host) map[string]struct{} {
	set := make(map[string]struct{})
	for ifidx, iface := range host.Ifaces {
		set[MACSetKey(ifidx, iface.MAC)] = struct{}{}
		set[MACSetKey(ifidx, BrdMAC)] = struct{}{}            // ip4 lookup
		set[MACSetKey(ifidx, SolMAC(iface.IP6))] = struct{}{} // ip6 lookup
	}
	return set
}
