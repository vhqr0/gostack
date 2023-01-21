package conf

import (
	"encoding/json"

	"github.com/vhqr0/gostack/lib/stack"
)

type StackConf struct {
	Host  *HostConf
	Route *RouteConf
}

func (conf *StackConf) NewStack() *stack.Stack {
	vhost := conf.Host.NewHost()
	vstack := stack.NewStack(vhost)
	conf.Route.AddRoute(vstack)
	return vstack
}

func (conf *StackConf) Marshal() ([]byte, error) {
	return json.Marshal(conf)
}

func (conf *StackConf) Unmarshal(data []byte) error {
	return json.Unmarshal(data, conf)
}
