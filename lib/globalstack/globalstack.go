package globalstack

import (
	"github.com/vhqr0/gostack/lib/conf"
	"github.com/vhqr0/gostack/lib/monitor"
	"github.com/vhqr0/gostack/lib/sock"
	"github.com/vhqr0/gostack/lib/stack"
)

var (
	Stack   *stack.Stack
	Monitor *monitor.Monitor
)

func Init(confFileName, httpListenAddr string) {
	Stack = conf.StackFromFile(confFileName)
	Monitor = &monitor.Monitor{
		Stack: Stack,
		Addr:  httpListenAddr,
	}
}

func Run() {
	Stack.Run()
	go Monitor.Serve()
}

func NewSock(family, typ uint32) (sock.Sock, error) {
	return Stack.NewSock(family, typ)
}
