package main

import (
	"flag"

	"github.com/vhqr0/gostack/lib/globalstack"
)

var (
	confFileName   = flag.String("c", "config.json", "config file name")
	httpListenAddr = flag.String("http", ":1080", "monitor http listen addr")
)

func main() {
	flag.Parse()

	globalstack.Init(*confFileName, *httpListenAddr)
	globalstack.Run()

	ch := make(chan struct{})
	<-ch
}
