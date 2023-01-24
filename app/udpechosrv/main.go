package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/vhqr0/gostack/lib/conf"
	"github.com/vhqr0/gostack/lib/monitor"
	"github.com/vhqr0/gostack/lib/sock"
)

var (
	confFileName   = flag.String("c", "config.json", "config file name")
	httpListenAddr = flag.String("http", ":1080", "monitor http listen addr")
	echoListenAddr = flag.String("echo", "0.0.0.0:7", "echo server listen addr")
)

func main() {
	flag.Parse()

	vstack := conf.StackFromFile(*confFileName)
	vstack.Run()

	go func() {
		m := &monitor.Monitor{Stack: vstack}
		log.Fatal(http.ListenAndServe(*httpListenAddr, m))
	}()

	addr, family, err := sock.ResolveAddr(*echoListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s, err := vstack.NewSock(family, sock.SockDgram)
	if err != nil {
		log.Fatal(err)
	}
	if addr, err := s.Bind(addr); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("bind %v", addr)
	}
	buf := [4096]byte{}
	for {
		n, addr, err := s.ReadFrom(buf[:])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("recv from %v: %s", addr, string(buf[:n]))
		if _, err := s.WriteTo(buf[:n], addr); err != nil {
			log.Fatal(err)
		}
	}
}