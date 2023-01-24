package conf

import (
	"encoding/json"
	"io"
	"log"
	"os"

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

func StackConfFrom(vstack *stack.Stack) *StackConf {
	conf := &StackConf{
		Host:  HostConfFrom(vstack.Host),
		Route: RouteConfFrom(vstack),
	}
	return conf
}

func StackFromFile(fileName string) *stack.Stack {
	file, err := os.Open(*&fileName)
	if err != nil {
		log.Panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}
	stackConf := &StackConf{}
	if err := stackConf.Unmarshal(data); err != nil {
		log.Panic(err)
	}
	vstack := stackConf.NewStack()
	return vstack
}
