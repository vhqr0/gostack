package main

import (
	"flag"
	"log"
	"time"

	"github.com/vhqr0/gostack/lib/globalstack"
	"github.com/vhqr0/gostack/lib/sock"
)

var (
	confFileName   = flag.String("c", "config.json", "config file name")
	httpListenAddr = flag.String("http", ":1080", "monitor http listen addr")
	echoListenAddr = flag.String("echo", "10.0.0.1:7", "echo server listen addr")
	message        = flag.String("m", "hello, world", "echo message")
	interval       = flag.Uint("i", 3, "echo interval")
)

func main() {
	flag.Parse()

	globalstack.Init(*confFileName, *httpListenAddr)
	globalstack.Run()

	addr, family, err := sock.ResolveAddr(*echoListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s, err := globalstack.NewSock(family, sock.SockDgram)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Connect(addr); err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			if _, err := s.Write([]byte(*message)); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	}()
	buf := [4096]byte{}
	for {
		if n, err := s.Read(buf[:]); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("recv %s", string(buf[:n]))
		}
	}
}
