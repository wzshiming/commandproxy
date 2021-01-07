package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/wzshiming/commandproxy"
)

var network = "tcp"

func init() {
	flag.StringVar(&network, "n", network, "network")
	flag.Parse()
}

func main() {
	targets := flag.Args()
	if len(targets) == 0 {
		log.Fatalln("not target")
		return
	}
	conn, err := net.Dial(network, targets[0])
	if err != nil {
		log.Fatalln(err)
		return
	}
	var buf1, buf2 [32 * 1024]byte
	err = commandproxy.Tunnel(context.Background(), commandproxy.Stdio, conn, buf1[:], buf2[:])
	if err != nil {
		log.Fatalln(err)
		return
	}
}
