package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/vhqr0/gostack/lib/conf"
	"github.com/vhqr0/gostack/lib/monitor"
)

var (
	confFileName   = flag.String("c", "config.json", "config file name")
	httpListenAddr = flag.String("http", ":1080", "monitor http listen addr")
)

func main() {
	flag.Parse()

	vstack := conf.StackFromFile(*confFileName)
	vstack.Run()

	m := &monitor.Monitor{Stack: vstack}
	log.Fatal(http.ListenAndServe(*httpListenAddr, m))
}
