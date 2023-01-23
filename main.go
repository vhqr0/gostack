package main

import (
	"flag"
	"io"
	"log"
	"os"

	"net/http"

	"github.com/vhqr0/gostack/lib/conf"
	"github.com/vhqr0/gostack/lib/monitor"

	_ "github.com/vhqr0/gostack/lib/util"
)

var (
	confFile   = flag.String("c", "config.json", "config file name")
	httpListen = flag.String("http", ":1080", "monitor http listen address")
)

func main() {
	flag.Parse()

	file, err := os.Open(*confFile)
	if err != nil {
		log.Fatal(err)
	}
	buf, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var stackConf = conf.StackConf{}
	stackConf.Unmarshal(buf)
	vstack := stackConf.NewStack()
	vstack.Run()

	m := &monitor.Monitor{Stack: vstack}
	log.Fatal(http.ListenAndServe(*httpListen, m))
}
