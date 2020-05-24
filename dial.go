package commandproxy

import (
	"context"
	"net"
)

type DialProxyCommand []string

func (p *DialProxyCommand) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
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
	cp := ProxyCommand(ctx, proxy[0], proxy[1:]...)
	return cp.Stdio()
}
