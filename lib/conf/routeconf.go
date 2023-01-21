package conf

import (
	"log"

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
