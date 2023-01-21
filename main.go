package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/vhqr0/gostack/lib/conf"

	_ "github.com/vhqr0/gostack/lib/util"
)

var (
	confFile = flag.String("c", "config.json", "config file")
)

func main() {
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
	ch := make(chan struct{})
	<-ch
}
