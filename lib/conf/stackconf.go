package conf

import (
	"encoding/json"

	"github.com/vhqr0/gostack/lib/stack"
)

type StackConf struct {
	Host  *HostConf
	Route *RouteConf
}

func (conf *StackConf) Marshal() ([]byte, error) {
	return json.Marshal(conf)
}

func (conf *StackConf) Unmarshal(data []byte) error {
	return json.Unmarshal(data, conf)
}

func (conf *StackConf) NewStack() *stack.Stack {
	vhost := conf.Host.NewHost()
	vstack := stack.NewStack(vhost)
	conf.Route.AddRoute(vstack)
	return vstack
}

func ConfFromStack(vstack *stack.Stack) *StackConf {
	conf := &StackConf{
		Host: ConfFromHost(vstack.Host),
		Route: ConfFromRouteTable(
			vstack,
			vstack.IPStack.Route4Table,
			vstack.IPStack.Route6Table),
	}
	return conf
}
