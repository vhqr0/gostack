package globalstack

import (
	"log"
	"net/http"

	"github.com/vhqr0/gostack/lib/conf"
	"github.com/vhqr0/gostack/lib/monitor"
	"github.com/vhqr0/gostack/lib/stack"
	"github.com/vhqr0/gostack/lib/sock"
)

var Stack *stack.Stack

func Run(confFileName string, httpListenAddr string) {
	Stack = conf.StackFromFile(*&confFileName)
	Stack.Run()

	m := &monitor.Monitor{Stack: Stack}
	log.Fatal(http.ListenAndServe(*&httpListenAddr, m))
}

func NewSock(family, typ uint32) (sock.Sock, error) {
	return Stack.NewSock(family, typ)
}
