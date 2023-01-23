package conf

import (
	"log"

	"github.com/vhqr0/gostack/lib/l3"
	"github.com/vhqr0/gostack/lib/stack"
)

type RouteConf struct {
	AutoRoute bool
	Entries   []*RouteEntryConf
}

type RouteEntryConf struct {
	Ver   int
	Iface string
	Peer  string
	Net   string
}

func (conf *RouteConf) AddRoute(vstack *stack.Stack) {
	if conf.AutoRoute {
		vstack.AutoRoute()
	}
	for _, entryConf := range conf.Entries {
		entryConf.AddRoute(vstack)
	}
}

func (conf *RouteEntryConf) AddRoute(vstack *stack.Stack) {
	if err := vstack.AddRoute(conf.Ver, conf.Iface, conf.Peer, conf.Net); err != nil {
		log.Panic(err)
	}
}

func ConfFromRouteTable(vstack *stack.Stack, route4Table, route6Table *l3.RouteTable) *RouteConf {
	conf := &RouteConf{AutoRoute: false}
	for _, entry := range route4Table.Entries {
		conf.Entries = append(conf.Entries, ConfFromRouteEntry(vstack, 4, entry))
	}
	for _, entry := range route6Table.Entries {
		conf.Entries = append(conf.Entries, ConfFromRouteEntry(vstack, 6, entry))
	}
	return conf
}

func ConfFromRouteEntry(vstack *stack.Stack, ver int, entry *l3.RouteEntry) *RouteEntryConf {
	var peer, net = "", ""
	if entry.Peer != nil {
		peer = entry.Peer.String()
	}
	if ones, bits := entry.Net.Mask.Size(); ones != bits { // gateway
		net = entry.Net.String()
	}
	iface := vstack.Host.Ifaces[entry.IfIdx].Name
	conf := &RouteEntryConf{
		Ver:   ver,
		Iface: iface,
		Peer:  peer,
		Net:   net,
	}
	return conf
}
