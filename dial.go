package commandproxy

import (
	"context"
	"net"
)

type DialProxyCommand []string

func (p *DialProxyCommand) Format(network string, address string) []string {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	} else if host == "" {
		host = "0.0.0.0"
	}
	m := map[byte]string{
		'n': network,
		'h': host,
		'p': port,
	}
	proxy := make([]string, len(*p))
	copy(proxy, *p)
	for i := range proxy {
		proxy[i] = ReplaceEscape(proxy[i], m)
	}
	return proxy
}

func (p *DialProxyCommand) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	proxy := p.Format(network, address)
	cp := ProxyCommand(ctx, proxy[0], proxy[1:]...)
	return cp.Stdio()
}
