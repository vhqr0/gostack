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

func RouteConfFrom(vstack *stack.Stack) *RouteConf {
	conf := &RouteConf{AutoRoute: false}
	table4, table6 := vstack.IPStack.Route4Table, vstack.IPStack.Route6Table
	for _, entry := range table4.Entries {
		conf.Entries = append(conf.Entries, RouteEntryConfFrom(vstack, 4, entry))
	}
	for _, entry := range table6.Entries {
		conf.Entries = append(conf.Entries, RouteEntryConfFrom(vstack, 6, entry))
	}
	return conf
}

func RouteEntryConfFrom(vstack *stack.Stack, ver int, entry *l3.RouteEntry) *RouteEntryConf {
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
