package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/google/tcpproxy"
)

var (
	listen  = flag.String("listen", "", "listen configuration, format: ':1234->foo.com,:2345->bar.com'")
	timeout = flag.Int("timeout", 10, "DNS cache timeout in seconds")
)

func main() {
	flag.Parse()

	if listen == nil || *listen == "" {
		log.Fatal("listen parameter must be set")
	}

	targets := strings.Split(*listen, ",")

	proxy := &tcpproxy.Proxy{}
	for _, target := range targets {
		targetParts := strings.Split(strings.TrimSpace(target), "->")
		if len(targetParts) != 2 {
			log.Fatalf("Invalid target: %s", target)
		}
		listen := strings.TrimSpace(targetParts[0])
		addr := strings.TrimSpace(targetParts[1])

		resolver := tcpproxy.NewDNSSRVResolver(
			addr,
			time.Duration(*timeout)*time.Second,
		)
		proxy.AddRoute(listen, &tcpproxy.DialProxy{AddrResolver: resolver})
	}
	log.Fatal(proxy.Run())
}
