package l3

import (
	"log"
	"net"

	"github.com/vhqr0/gostack/lib/host"
)

type RouteTable struct {
	Entries []*RouteEntry
}

type RouteEntry struct {
	IfIdx int
	Peer  net.IP
	Net   *net.IPNet
}

func (table *RouteTable) Next(peer net.IP) *RouteEntry {
	for _, entry := range table.Entries {
		if entry.Net.Contains(peer) {
			return entry
		}
	}
	return nil
}

func (table *RouteTable) Add(entry *RouteEntry) {
	table.Entries = append(table.Entries, entry)
}

func (stack *IPStack) Next(peer net.IP) (entry *RouteEntry) {
	switch len(peer) {
	case 4:
		entry = stack.Route4Table.Next(peer)
	case 16:
		entry = stack.Route6Table.Next(peer)
	default:
		log.Panic(&InvalidIPLen{Len: len(peer)})
	}
	return
}

func (stack *IPStack) AutoRoute() {
	for ifidx, iface := range stack.EthStack.Host.Ifaces {
		entry4 := &RouteEntry{IfIdx: ifidx, Net: iface.Net4}
		entry6 := &RouteEntry{IfIdx: ifidx, Net: iface.Net6}
		stack.Route4Table.Add(entry4)
		stack.Route6Table.Add(entry6)
	}
}

func (stack *IPStack) AddRoute(ver int, ifname, peerStr, netStr string) error {
	var (
		err   error
		peer  net.IP
		ipnet *net.IPNet
	)

	ifidx := stack.EthStack.Host.GetIfaceByName(ifname)

	if ifidx < 0 {
		return &host.InvalidIfaceName{Name: ifname}
	}

	if peerStr != "" { // have peer
		if peer, err = ParseIP(ver, peerStr); err != nil {
			return err
		}
	}

	if netStr != "" { // have net
		if _, ipnet, err = ParseCIDR(ver, netStr); err != nil {
			return err
		}
	} else { // use all ip net (gateway)
		if ipnet, err = AllIPNet(ver); err != nil {
			return err
		}
	}

	entry := &RouteEntry{
		IfIdx: ifidx,
		Peer:  peer,
		Net:   ipnet,
	}

	switch ver {
	case 4:
		stack.Route4Table.Add(entry)
	case 6:
		stack.Route6Table.Add(entry)
	default:
		return &InvalidIPVer{Ver: ver}
	}

	return nil
}
