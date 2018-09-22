package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/tcpproxy"
)

var (
	listen    = flag.String("listen", "", "listen configuration, format: ':1234->foo.com,:2345->bar.com'")
	timeout   = flag.Int("timeout", 10, "DNS cache timeout in seconds")
	dnsServer = flag.String("dns-server", "", "DNS server to use")
)

func main() {
	flag.Parse()

	if listen == nil || *listen == "" {
		log.Fatal("listen parameter must be set")
	}

	targets := strings.Split(*listen, ",")

	dnsResolver := &net.Resolver{}
	if dnsServer != nil && *dnsServer != "" {
		dnsResolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			d := &net.Dialer{}
			return d.DialContext(ctx, "udp", "8.8.8.8:53")
		}
	}

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
			dnsResolver,
		)
		proxy.AddRoute(listen, &tcpproxy.DialProxy{AddrResolver: resolver})
	}
	log.Fatal(proxy.Run())
}
